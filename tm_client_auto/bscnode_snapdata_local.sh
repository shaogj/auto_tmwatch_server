#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


TRANS_DATA_TO="$1"
DESC_IP="$2"
DESC_SNAP_DIR="$3"
bscgethdataset=`/data/node/geth`

function printHelp () {
	echo "Usage: ./bscnode_snapdata <from|to> <getsnap-descIP> <snaptimedir>.\nThe arguments must be in order."
}

function validateArgs () {
	if [ -z "${TRANS_DATA_TO}" ]; then
		echo "Option from / to / snaptimedir / trans snapdata to bsc node"
		printHelp
		exit 1
	fi
	if [ -z "${DESC_IP}" ]; then
		echo "setting to default DESC_IP '192.169.1.333'"
		DESC_IP=mychannel
	fi
}

function restartBSCNode () {
  	echo "to restart local BSCNode:"
     docker restart trust_zkbsc
     #add time sleep
     sleep 30
}

#停止将要拷贝的正常运行的bsc节点(数据源节点)
function snapData() {
    sudo docker stop trust_zkbsc
}

#恢复拷贝的快照数据
function RecoverNormalDataFromSnap() {
      echo "step2,cur exec RecoverDataFromSnap proc.RecoverDataFromSnap is: $curBscSnapTime"
  	  afterSnapData="cur time is:$(date "+%Y%m%d%H%M%S"),get snapdata from:${DESC_IP},finished,to replace data to local bsc gethdata!"
      echo $afterSnapData
      #std--
      bsc_gethdata="/data/node/geth"
      #0425testing
      bsc_gethdata="/Users/gejianspro/go/src/202108FromBFLProj/ChainWatch_Project2023"
      #0417localing
      #bsc_gethdata="/Users/gejians/go/src/2021New_BFLProjTotal/0424NewTMEnvRes/20230417scp_newgeth"
      echo $bsc_gethdata
      curBscSnapTimebak=$(date +%Y%m%d%H%M)
      #real data
      mv $bsc_gethdata/chaindata $bsc_gethdata/chaindata_baklast_${curBscSnapTimebak}
      mv $bsc_gethdata/triecache $bsc_gethdata/triecache_baklast_${curBscSnapTimebak}
      echo "get copy RecoverData,FromSnap dir: $curBscSnapTime"
      cp -r ./chaindata_${DESC_SNAP_DIR}.tar.gz ${bsc_gethdata}
      cp -r ./triecache_${DESC_SNAP_DIR}.tar.gz ${bsc_gethdata}
      echo "get RecoverDataFromSnap finished!"
      tar -zxvf $bsc_gethdata/chaindata_${DESC_SNAP_DIR}.tar.gz -C $bsc_gethdata
      tar -zxvf $bsc_gethdata/triecache_${DESC_SNAP_DIR}.tar.gz -C $bsc_gethdata
      echo "get tar file to datadir finished!"


}

#online
function GetSnapDataFromRemoteNode() {
 #从远程拷贝数据
  echo "to exec GetSnapDataFromRemoteNode(),to scp bsc snapdate from node，IP is:"
  echo $DESC_IP
  echo "to scp from desc ip dir is:${DESC_SNAP_DIR}"
  echo $DESC_SNAP_DIR

  echo "get curBscSnapTime value is: $curBscSnapTime"
  mkdir $curBscSnapTime
  #cd /data/node/geth
  cd $curBscSnapTime
  echo "cur work dir curBscSnapTime value is: ${curBscSnapTime},start scp datafrom remote ip!"
  #real data
# scp -r root@${DESC_IP}:/data/node/geth/chaindata ./
# scp -r root@${DESC_IP}:/data/node/geth/triecache ./
  #real data tardata
  #0425doing
  scp -r root@${DESC_IP}:/data/node/geth/chaindata_${DESC_SNAP_DIR}.tar.gz ./
  scp -i /root/bscWatchTest/Key_dev-user_rsa-Local-0422 -r root@${DESC_IP}:/data/node/geth/triecache_${DESC_SNAP_DIR}.tar.gz ./
   #180node,testdata
  #scp -i /home/dev-user/BJdev-user_rsa-0628 -r dev-user@${DESC_IP}:/data/node/geth/nodes ./
  #scp -i /home/dev-user/BJdev-user_rsa-0628 -r dev-user@${DESC_IP}:/data/node/geth/check_snapcp04 ./
  afterSnapData="cur work dir:${curBscSnapTime},get snapdata from:${DESC_IP} finished!"
  echo $afterSnapData

}
validateArgs

#Create the network using docker compose
if [ "${TRANS_DATA_TO}" == "to" ]; then
	echo "to exec scp bsc snapdate to Desc IP node"
	echo $DESC_IP


elif [ "${TRANS_DATA_TO}" == "from" ]; then ## Clear the network
	echo "to exec scp to get bsc snapdata from IP node: $DESC_IP"

	curBscSnapTime=$(date +%Y%m%d%H%M)
	echo "cur scpsnaptime is: $curBscSnapTime"
	echo "get invoke request' params num is :$#"
	#0410add
    if [ $# != 3 ] ; then
      echo "USAGE: $0 from to--0411testing===params num is no 3!"
      printHelp
      exit 1
    fi
   if [ "${DESC_SNAP_DIR}" == "default" ]; then ## Clear the network
   #,DESC_SNAP_DIR
   	 echo "get DESC_SNAP_DIR' params num is :$DESC_SNAP_DIR"
   fi
	echo "get DESC_SNAP_DIR' request params num is :$DESC_SNAP_DIR"
	#exit 1
  GetSnapDataFromRemoteNode
  #0417,exit 1
	RecoverNormalDataFromSnap
	echo "to exec docker restart trust_zkbsc,curtime is: $curBscSnapTime"
	docker restart trust_zkbsc
	curBscSnapTimeafter=$(date +%Y%m%d%H%M)

	echo "after exec docker restart trust_zkbsc,curtime is: curBscSnapTimeafter"

else
	printHelp
	exit 1
fi

