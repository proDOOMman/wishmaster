#!/bin/sh

# Source function library
. /etc/init.d/init-functions

DAEMON=/mnt/persistent/wishmaster/wishmaster
BASENAME=wishmaster
DESC="wishmaster"

ARGS=""

# Each init script action is defined below...
start() {
	local RET

	log_status_msg "Starting $DESC $ARGS : " -n

	start-stop-daemon -S -b --quiet --exec env "LD_PRELOAD=/lib/libpthread.so.0" "ADDR=0.0.0.0" "PORT=80" $DAEMON -- $ARGS
	RET=$?
	if [ $RET -ne 0 ]; then
		log_failure_msg " [failed]"
		return 1
	fi

 	log_success_msg " [done]"

	return 0
}

stop() {
	local RET

	log_status_msg "Stopping $DESC: " -n
  	start-stop-daemon -K --quiet --name $BASENAME

  	RET=$?
	if [ $RET -ne 0 ]; then
		log_failure_msg " [failed]"
		return 1
	fi

  	log_success_msg " [done]"
	return 0
}

restart() {
	local RET

	log_status_msg "Restarting $DESC..."
	stop
	sleep 2
	start
	RET=$?

	return $RET
}


case "$1" in
	start)
		start
		return $?
		;;
	stop)
		stop
		return $?
		;;
	restart)
		restart
		return $?
		;;

       *)
    	        echo "Usage: $0 {start|stop|restart}"
		;;
esac

return 1
