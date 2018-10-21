#/bin/bash

_dot_completions() {
#	COMPREPLY=($(${COMP_WORDS[0]} __completebash -c ${COMP_CWORD} "${COMP_WORDS[@]}"))
#	COMPREPLY="line ${COMP_LINE} point ${COMP_POINT}"
#	cur=${COMP_WORDS[$COMP_CWORD]}
#	if [ "$cur" == "a" ]; then
#		COMPREPLY="$cur"
#		compopt -o nospace
#	else
#		COMPREPLY="$cur"
#	fi
#	_longopt
	local base=$1 cur=$2 prev=$3
#	COMPREPLY="base[$base] cur[$cur] prev[$prev]"
#	compopt -o nosort
	mapfile -t COMPREPLY < <(grep "^$cur" <<EOF
first word
second longer phrase
and last line here
EOF
)
	if [ ${#COMPREPLY[@]} -eq 1 ]; then
		a="${COMPREPLY[0]}"
		[[ "$a" =~ " " ]] && COMPREPLY[0]=`printf '"%s"' "$a"`
	fi
}

complete -F _dot_completions dot
