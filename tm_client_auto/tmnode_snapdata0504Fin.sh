#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


TRANS_DATA_TO="$1"
DESC_SNAP_DIR="$2"
SNAP_DATA_DIR="$3"
bscgethdataset=`/data/node/geth`

function printHelp () {
	echo "Usage: ./tmnode_snapdata <restoredata|baksnapdata>\nThe arguments must be in order."
}

function validateArgs () {
	if [ -z "${TRANS_DATA_TO}" ]; then
		echo "Option restoredata / baksnapdata /trans snapdata to tm node"
		printHelp
		exit 1
	fi
	if [ -z "${SNAP_DATA_DIR}" ]; then
		echo "setting to default snapdata '/bakdata/data'"
		SNAP_DATA_DIR='/bakdata/data'
	fi
}
function clearUselesstmfiles () {
  rm /trias/.ethermint/tendermint/config/write-file-atomic*
  rm /trias/.ethermint/tendermint/tm-logs/tmware.log
  rm /trias/.ethermint/tendermint/config/addrbook.json
}

function stoptmProc () {
  cd /trias/.ethermint/tendermint/
  echo "to stoptmProc tmNode,cur time is: $curBscSnapTime" >> tmrestoresnap.log
  timeout -k 3s 10s sudo systemctl stop tendermint
  tmexist=`ps aux|grep tendermint|grep -v grep`
  if [ -z $tmexist ];then
    echo "after systemctl stop tendermint,,tm process exist no!" >> tmrestoresnap.log
  else
    echo "tendermint  process still exist!,,to kill tm proc!"   >> tmrestoresnap.log
    kill -9 `ps aux|grep tendermint|grep -v grep|awk '{print $2}'`
    echo "cur kill proc tendermint----done!"    >> tmrestoresnap.log
    sleep 5
  fi
  tmexist=`ps aux|grep tendermint|grep -v grep`
  echo "after stop tenermint proc! tendermint pidinfo is :$tmexist"   >> tmrestoresnap.log
}

#停止将要拷贝的正常运行的tm节点,及监控服务(数据源节点)
function stopTmWatch () {
  echo "to stop tmwatch:"
  systemctl stop tmwatch
  echo "cur tmwatch is stopped"
  echo "to stop tendermint:"
  #sudo systemctl stop tendermint
  stoptmProc
  clearUselesstmfiles
  sleep 3
}

#启动tm节点,及监控服务
function restartTmSever () {
  	echo "to restart tendermint:"
  	sudo systemctl restart tendermint
  	sleep 10
   	echo "restart local tmNode finished!,to restart tmwatch proc"
    sudo systemctl restart tmwatch
     #add time sleep
}

#online
#从本地的tm快照备份，拷贝数据
function GetSnapDataFromRemoteNode() {
  cd /trias/.ethermint/tendermint/
  echo "to exec GetSnapDataFromRemoteNode(),to get tm snapdate from node，bakdata from dir is: $SNAP_DATA_DIR " >> tmrestoresnap.log
  echo "to scp from desc tmdata dir is:${DESC_SNAP_DIR}" >> tmrestoresnap.log
  echo $DESC_SNAP_DIR

  echo "get curBscSnapTime value is: $curBscSnapTime" >> tmrestoresnap.log

  #cd $curBscSnapTime
  echo "cur work dir curBscSnapTime value is: ${curBscSnapTime},start restore tm snapdata from localbak!"
  #0419tem,,0420:
  echo "record cur tmheight info of priv_validator_state.json is:" >> tmrestoresnap.log
  cat /trias/.ethermint/tendermint/data/priv_validator_state.json >> tmrestoresnap.log
  mv /trias/.ethermint/tendermint/data /trias/.ethermint/tendermint/data_${curBscSnapTime}_bak
  #mkdir /trias/.ethermint/tendermint/snaptmp
  #cp -r ${SNAP_DATA_DIR}/${DESC_SNAP_DIR}/data /trias/.ethermint/tendermint/snaptmp
  cp -r ${SNAP_DATA_DIR}/${DESC_SNAP_DIR}/data /trias/.ethermint/tendermint/
  #检查cur height after recover:
  echo "record cur restore snapdata info of priv_validator_state.json,in ${DESC_SNAP_DIR} is:" >> tmrestoresnap.log
  cat /trias/.ethermint/tendermint/data/priv_validator_state.json >> tmrestoresnap.log

  afterSnapData="cur work dir:${curBscSnapTime},get snapdata from:${DESC_SNAP_DIR} finished!"
  echo $afterSnapData >> tmrestoresnap.log

}
validateArgs

#Create the network using docker compose
if [ "${TRANS_DATA_TO}" == "to" ]; then
	echo "to exec get tm baksnapdate to cur tmdata"
	echo $SNAP_DATA_DIR


elif [ "${TRANS_DATA_TO}" == "restoredata" ]; then ## Clear the network
	echo "to exec scp to get tm snapdata from tmdatadir: $DESC_SNAP_DIR"

	curBscSnapTime=$(date +%Y%m%d%H%M)
	echo "cur scpsnaptime is: $curBscSnapTime"
	echo "get invoke request' params num is :$#"
	#0410add,if [ $# != 3 ] ; then
	#if[ $x -le5]
   if [ $# != 2 ] ; then
      echo "USAGE: $0 from to--0411testing===params num less then 2!"
      printHelp
      exit 1
   fi
	echo "get DESC_SNAP_DIR' request params num is :$DESC_SNAP_DIR"

	stopTmWatch
  #0419stem,check..
  GetSnapDataFromRemoteNode
  #to do:
  #restartTmSever
	echo "to exec restart tendermint,curtime is: $curBscSnapTime"
	curBscSnapTimeafter=$(date +%Y%m%d%H%M)
	#echo "after exec docker restart trust_zkbsc,curtime is: $curBscSnapTimeafter"

else
	printHelp
	exit 1
fi

