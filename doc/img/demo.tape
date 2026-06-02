# omokage demo tape (https://github.com/charmbracelet/vhs)
#
# Regenerate with `make demo`, which builds the CLI, seeds a throwaway project
# under /tmp/omokage-demo, runs this tape, and moves demo.gif into doc/img/.
#
# vhs v0.11 parses absolute Output paths oddly, so the tape writes a relative
# demo.gif into the current directory and `make demo` moves it into doc/img/.

Output demo.gif

Set Shell bash
Set FontSize 18
Set Width 1180
Set Height 820
Set Padding 24
Set Theme "Catppuccin Mocha"

Hide
Type "export PATH=/tmp:$PATH; cd /tmp/omokage-demo; clear" Enter
Show

Type "omokage init" Enter
Sleep 1200ms

Type "omokage train --author me examples/posts" Enter
Sleep 1500ms

Type "# A new draft that keeps the author's voice" Enter
Type "omokage check --author me examples/draft-keeps-voice.md" Enter
Sleep 2500ms

Type "# The same idea rewritten in a stiff, different voice" Enter
Type "omokage check --author me examples/draft-lost-voice.md" Enter
Sleep 2500ms

Type "omokage diff examples/draft-keeps-voice.md examples/draft-lost-voice.md" Enter
Sleep 2800ms
