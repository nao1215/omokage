#!/bin/sh
# Writes the inputs of the himorime suite in bench/. Deterministic: the same
# arguments always write the same bytes, so a base and a head revision read
# the same input.
#
#   sh gen.sh doc P SALT FILE     one post of P mixed Japanese and English paragraphs
#   sh gen.sh corpus N P DIR      N posts of P paragraphs, each slightly different
set -eu

doc() {
	awk -v p="$1" -v salt="$2" 'BEGIN {
		for (i = 0; i < p; i++) {
			printf "私は朝のしずかな時間が好きです。白湯を一杯ゆっくり飲みながら、窓の外をながめます。"
			printf "On quiet mornings I drink a cup of warm water and watch the light change outside. "
			printf "特別なことは何もしていません。それでも、こうした小さな習慣が一日をおだやかにする気がしています。\n\n"
		}
		printf "追記その%d。\n", salt
	}' > "$3"
}

case "$1" in
doc)
	mkdir -p "$(dirname "$4")"
	doc "$2" "$3" "$4"
	;;
corpus)
	mkdir -p "$4"
	i=0
	while [ "$i" -lt "$2" ]; do
		doc "$3" "$i" "$4/post$(printf '%03d' "$i").md"
		i=$((i + 1))
	done
	;;
*)
	echo "gen.sh: unknown kind $1" >&2
	exit 2
	;;
esac
