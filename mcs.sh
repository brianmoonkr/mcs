#!/bin/sh

PID_FILE='PID'
COMMAND=$1

USER_HOME='/home/teamgrit'
RUN_FILE="${USER_HOME}/Projects/mcs/mcs"

#echo $! > $PID_FILE

if [ -f ${PID_FILE} ] ; then
	exec < $PID_FILE
	read PID
	if [ ! "x${PID}" = "x" ] && kill -0 ${PID} 2>/dev/null; then
		STATUS="(pid ${PID}) already running"
		RUNNING=1
	else
		STATUS="(pid ${PID}) not running"
		RUNNING=0
	fi	
else
	STATUS="(no pid file) not running"
	RUNNING=0
	PID=-
fi


ERROR=0
case ${COMMAND} in
'start')
	if [ ${RUNNING} -eq 1 ]; then
		echo "$0 ${COMMAND}: ${STATUS}"
		ERROR=2
	else
		${RUN_FILE} 1>&2 &
		if kill -0 $! 2>/dev/null; then
			echo $! > $PID_FILE
			echo "$0: started. PID:$!"
		else
			echo "$0: could not be started."
			ERROR=3
		fi
	fi
	;;

'stop')
	if [ ${RUNNING} -eq 0 ]; then
		echo "$0 ${COMMAND}: ${STATUS}"
		rm ${PIDFILE} 1>/dev/null 2>&1
		ERROR=4
	else 
		if kill -9 ${PID} 2>/dev/null; then
			echo "$0 ${COMMAND}: stopped"
			rm ${PID_FILE} 1>/dev/null 2>&1
		else
			echo "$0 ${COMMAND}: could not be stopped"
			ERROR=5
		fi
	fi
	;;

'restart')
	if [ ${RUNNING} -eq 0 ]; then
		echo "$0 ${COMMAND}: not running, trying to start"
		${RUN_FILE} 1>&2 &
		if kill -0 $! 2>/dev/null; then
			echo $! > $PID_FILE
			echo "$0: started. PID:$!"
		else
			echo "$0: could not be started."
			ERROR=6
		fi
	else
		${RUN_FILE} 1>&2 &
		if kill -0 $! 2>/dev/null; then
			echo $! > ${PID_FILE}
			if kill -KILL ${PID} 2>/dev/null; then
				echo "$0 ${COMMAND}: restarted. PID:$!"
			else
				echo "$0 ${COMMAND}: could not be restarted."
				ERROR=7
			fi
		else
			echo "$0 ${COMMAND}: could not be restarted."
			ERROR=8
		fi
	fi
	;;

'status')
	echo ${STATUS}
	;;

'check')
	echo ${PID} ${RUNNING}
	;;

*)
	echo "usage: $0 (start|stop|restart|status|check)"

esac


exit ${ERROR}

