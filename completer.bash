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
	echo "$@"
	echo qwe
}

complete -F _dot_completions dot
