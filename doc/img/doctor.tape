# omokage doctor demo tape (https://github.com/charmbracelet/vhs)
#
# Regenerate with `make demo`, which builds the CLI, seeds a throwaway project
# under /tmp/omokage-demo, runs this tape, and moves doctor.gif into doc/img/.
# The GIF focuses on the corpus-quality check: is a corpus solid enough to train,
# and if not, what to do about it. vhs v0.11 parses absolute Output paths oddly,
# so the tape writes a relative doctor.gif and `make demo` moves it into doc/img/.

Output doctor.gif

Set Shell bash
Set FontSize 16
Set Width 1200
Set Height 900
Set Padding 24
Set Theme "Catppuccin Mocha"

Hide
Type "export PATH=/tmp:$PATH; cd /tmp/omokage-demo; clear" Enter
Show

Type "# doctor reads a corpus and rates how reliable its scores will be — it trains nothing" Enter
Type "omokage doctor examples/ja/posts" Enter
Sleep 3500ms

Type "# Too few samples? doctor says so, and what to do next" Enter
Type "omokage doctor examples/ja/draft-keeps-voice.md examples/ja/draft-lost-voice.md" Enter
Sleep 3800ms
