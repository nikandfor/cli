
# %[1]s bash completion
_nikandcli_complete_bash() {
	declare -x CLI_COMP_BASE=$1
	declare -x CLI_COMP_CUR=$2
	declare -x CLI_COMP_PREV=$3

	declare -x CLI_COMP_LINE=$COMP_LINE
	declare -x CLI_COMP_INDEX=$COMP_POINT

	declare -x CLI_COMP_WORDS_LENGTH=${#COMP_WORDS[@]}
	declare -x CLI_COMP_WORDS_INDEX=$COMP_CWORD

	for i in $(seq 0 $(("${#COMP_WORDS[@]}" - 1)) ); do
		declare -x "CLI_COMP_WORDS_${i}"="${COMP_WORDS[$i]}"
	done

	cmd=$("$1" "${COMP_WORDS[@]:1:$COMP_CWORD}")

	eval "$cmd"
}

complete -F _nikandcli_complete_bash greeter # %[2]s

#_nikandcli_complete_bash "$@"

# to persist bash completion add this to the end of your ~/.bashrc file by command:
#   %[2]s >>~/.bashrc
# or alternatively to enable it to only current session use command:
#   source <(%[2]s)
