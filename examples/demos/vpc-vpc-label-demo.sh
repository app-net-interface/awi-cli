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
DEMO_CMD_COLOR=$BLACK
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

Demo Steps
==========

1. Devops admin creates two VPCs (staging and development)
with three instances (two databases and dashboard) on AWS using Terraform.

2. Checks connectivity between VPCs – not connected.

3. Devops admin uses AWI CLI (with vManage Plugin) to request connectivity to both VPCs with 'deny all' policy.

4. vManage provisions connectivity on behalf of the devops admin.

5. Checks connectivity between VPCs – not reachable due to 'deny all' policy.

6. Devops admin uses AWI CLI (with vManage Plugin) to request connectivity between database apps.

7. After connectivity is provisioned, staging databases can reach development database instances but not the other one.

EOM

read -p ""
echo "For this demo we have pre-created an AWS VPCs and an EC2 instances using terraform."
read -p ""
clear
echo "Showing the AWS VPC Ids"
read -p ""
pe "./awi list vpc-tag --cloud aws"
echo "Showing the AWS VPC Instances used in demo"
read -p ""
pe "./awi list instance --cloud aws --tag 'staging*'"
pe "./awi list instance --cloud aws --tag 'development*'"
echo "Showing current connection (none)"
read -p ""
pe "./awi list connection"
echo "We will try to reach development VM from staging VM."
read -p ""
echo "We can't connect to instance."
read -p "Now we will request connectivity using vManage CLI using AWI interface."
echo "Before that, let's see the AWI Config that we will use to request this connectivity"
read -p ""
pe "cat vpc-vpc.yaml"
read -p ""
clear
read -p "Running AWI CLI ..."
pe "./awi connect --connection-config vpc-vpc.yaml"
pe "./awi list connection"
echo "We will try to reach the instance with requested label (database)."
read -p ""
echo "We can't connect to instance due to 'deny' policy."
read -p "We will request connectivity for database app."
pe "cat vpc-vpc-label-acl.yaml"
read -p ""
clear
read -p "Running AWI CLI ..."
pe "./awi connect --connection-config vpc-vpc-label-acl.yaml"
echo "We will try again to reach the instance with requested label (database)."
pe "./awi list instance --cloud aws --tag 'database*'"
read -p ""
echo "Connections have been established."
echo "Now we will try to reach the instance without requested label (dashboard)."
pe "./awi list instance --cloud aws --tag 'dashboard*'"
read -p ""
echo "We can't connect to instance due to 'deny' policy."
read -p ""
