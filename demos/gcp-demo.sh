#!/bin/bash

###############################################################################
#
# demo-magic.sh
#
# Copyright (c) 2015 Paxton Hare
#
# This script lets you script demos in bash. It runs through your demo script when you press
# ENTER. It simulates typing and runs commands.
#
###############################################################################

# the speed to "type" the text
TYPE_SPEED=20

# no wait after "p" or "pe"
NO_WAIT=true

# if > 0, will pause for this amount of seconds before automatically proceeding with any p or pe
PROMPT_TIMEOUT=0

# don't show command number unless user specifies it
SHOW_CMD_NUMS=false


# handy color vars for pretty prompts
BLACK="\033[0;30m"
BLUE="\033[0;34m"
GREEN="\033[0;32m"
GREY="\033[0;90m"
CYAN="\033[0;36m"
RED="\033[0;31m"
PURPLE="\033[0;35m"
BROWN="\033[0;33m"
WHITE="\033[1;37m"
COLOR_RESET="\033[0m"

C_NUM=0

# prompt and command color which can be overridden
DEMO_PROMPT="$ "
DEMO_CMD_COLOR=$WHITE
DEMO_COMMENT_COLOR=$GREY

##
# prints the script usage
##
function usage() {
  echo -e ""
  echo -e "Usage: $0 [options]"
  echo -e ""
  echo -e "\tWhere options is one or more of:"
  echo -e "\t-h\tPrints Help text"
  echo -e "\t-d\tDebug mode. Disables simulated typing"
  echo -e "\t-n\tNo wait"
  echo -e "\t-w\tWaits max the given amount of seconds before proceeding with demo (e.g. '-w5')"
  echo -e ""
}

##
# wait for user to press ENTER
# if $PROMPT_TIMEOUT > 0 this will be used as the max time for proceeding automatically
##
function wait() {
  if [[ "$PROMPT_TIMEOUT" == "0" ]]; then
    read -rs
  else
    read -rst "$PROMPT_TIMEOUT"
  fi
}

##
# print command only. Useful for when you want to pretend to run a command
#
# takes 1 param - the string command to print
#
# usage: p "ls -l"
#
##
function p() {
  if [[ ${1:0:1} == "#" ]]; then
    cmd=$DEMO_COMMENT_COLOR$1$COLOR_RESET
  else
    cmd=$DEMO_CMD_COLOR$1$COLOR_RESET
  fi

  # render the prompt
  x=$(PS1="$DEMO_PROMPT" "$BASH" --norc -i </dev/null 2>&1 | sed -n '${s/^\(.*\)exit$/\1/p;}')

  # show command number is selected
  if $SHOW_CMD_NUMS; then
   printf "[$((++C_NUM))] $x"
  else
   printf "$x"
  fi

  # wait for the user to press a key before typing the command
  if [ $NO_WAIT = false ]; then
    wait
  fi

  if [[ -z $TYPE_SPEED ]]; then
    echo -en "$cmd"
  else
    echo -en "$cmd" | pv -qL $[$TYPE_SPEED+(-2 + RANDOM%5)];
  fi

  # wait for the user to press a key before moving on
  if [ $NO_WAIT = false ]; then
    wait
  fi
  echo ""
}

##
# Prints and executes a command
#
# takes 1 parameter - the string command to run
#
# usage: pe "ls -l"
#
##
function pe() {
  # print the command
  p "$@"
  run_cmd "$@"
}

##
# print and executes a command immediately
#
# takes 1 parameter - the string command to run
#
# usage: pei "ls -l"
#
##
function pei {
  NO_WAIT=true pe "$@"
}

##
# Enters script into interactive mode
#
# and allows newly typed commands to be executed within the script
#
# usage : cmd
#
##
function cmd() {
  # render the prompt
  x=$(PS1="$DEMO_PROMPT" "$BASH" --norc -i </dev/null 2>&1 | sed -n '${s/^\(.*\)exit$/\1/p;}')
  printf "$x\033[0m"
  read command
  run_cmd "${command}"
}

function run_cmd() {
  function handle_cancel() {
    printf ""
  }

  trap handle_cancel SIGINT
  stty -echoctl
  eval "$@"
  stty echoctl
  trap - SIGINT
}


function check_pv() {
  command -v pv >/dev/null 2>&1 || {

    echo ""
    echo -e "${RED}##############################################################"
    echo "# HOLD IT!! I require pv but it's not installed.  Aborting." >&2;
    echo -e "${RED}##############################################################"
    echo ""
    echo -e "${COLOR_RESET}Installing pv:"
    echo ""
    echo -e "${BLUE}Mac:${COLOR_RESET} $ brew install pv"
    echo ""
    echo -e "${BLUE}Other:${COLOR_RESET} http://www.ivarch.com/programs/pv.shtml"
    echo -e "${COLOR_RESET}"
    exit 1;
  }
}

check_pv
#
# handle some default params
# -h for help
# -d for disabling simulated typing
#
while getopts ":dhncw:" opt; do
  case $opt in
    d)
      unset TYPE_SPEED
      ;;
    n)
      NO_WAIT=true
      ;;
    c)
      SHOW_CMD_NUMS=true
      ;;
    w)
      PROMPT_TIMEOUT=$OPTARG
      ;;
    *)
      usage
      exit 1
      ;;
  esac
done

clear

cat << EOM

In this demo we will show how to setup connectivity between Google Cloud Platform VPCs and applications running in those VPCs
using MyCelium: Application WAN Interface system.

EOM

read -p ""
read -p "In our setup we have couple of network domains."
pe "myCelium list network-domains"
read -p ""
read -p "For now we don't have any connections between network domains created."
pe "myCelium list connection"
read -p ""
clear
read -p "In our demo we will try to connect two GCP VPCs: staging VPC and development VPC."
pe "myCelium list vpc --cloud gcp"
read -p ""
read -p "In both staging and development VPCs we have 3 instances: one 'dashboard' and two 'databases'."
pe "myCelium list instance --cloud gcp"
read -p ""
read -p "We will try to ping development VMs from staging first."
read -p "Ping failed, because there's no connection between those domains."
clear
read -p "Now we will request connectivity between staging and development VPCs with myCelium CLI using AWI interface."
read -p "Let's see the AWI Config that we will use to request this connectivity."
read -p "We will set default access to allow all traffic."
pe "cat examples/gcp-vpc-vpc-demo-allow.yaml"
read -p ""
read -p "Now let's run myCelium CLI to create connection..."
pe "myCelium create connection --connection-config examples/gcp-vpc-vpc-demo-allow.yaml"
read -p ""
read -p "Establishing connection takes couple of minutes during which recording will be paused..."
read -p "Let's list connections to see if our connection is applied already."
pe "myCelium list connection"
read -p ""
read -p "Connection is bidirectional so reverse connection from development VPC to staging VPC was automatically created by Cisco SD-WAN controller."
read -p "Now let's check if there's connectivity between staging and development."
read -p "As expected we can ping all VMs from development VPC."
read -p "Now we will disconnect staging and development VPCs..."
pe "myCelium delete connection 8528197605643627928:1747121394098482639"
read -p ""
pe "myCelium list connection"
read -p ""
read -p "Ping is not working now"
clear
read -p "Now we will request connectivity between staging and development VPCs again, but this time with default access set to deny."
read -p "Let's see the AWI Config that we will use to request this connectivity."
pe "cat examples/gcp-vpc-vpc-demo-deny.yaml"
read -p ""
read -p "Now let's run AWI CLI to create connection..."
pe "myCelium create connection --connection-config examples/gcp-vpc-vpc-demo-deny.yaml"
read -p ""
read -p "And let's list connections to see if it's applied."
pe "myCelium list connection"
read -p ""
read -p "Now let's check if there's connectivity between staging and development."
read -p "As expected this time ping failed, all traffic is denied by default."
clear
read -p "Now we will request application connectivity only between staging and development database instances."
read -p "Let's see the AWI Config that we will use to request this connectivity."
pe "cat examples/gcp-vpc-vpc-label-app-connection.yaml"
read -p ""
read -p "Connectivity will be established based on labels, instances which have matching labels will be connected."
read -p "Now let's run AWI CLI to create connection..."
pe "myCelium create app-connection --connection-config examples/gcp-vpc-vpc-label-app-connection.yaml"
read -p ""
read -p "And let's list app connections to see if it's applied."
pe "myCelium list app-connection"
read -p ""
read -p "Now let's check if there's connectivity between staging and development databases.
Connections between dashboard instances and dashboards to databases should be discarded."
clear
read -p "Now we will create new VM in staging VPC, it will be new database instance.
Application connectivity should be automatically enabled, we should be able to ping development database instances, and not able to ping development dashboard."
clear
read -p "Now we will again create new VM in staging VPC, this time it will be new dashboard instance.
Application connectivity should be automatically enabled, we should NOT be able to ping any development VM."
read -p "When we will create destination VM app connectivity rules will also be automatically added."
read -p "After removing app-connection connectivity will be disabled."
clear
read -p "In this demo we showed how to create connection with GCP VPCs and instances using Application WAN Interface.
Thank you for watching!"
