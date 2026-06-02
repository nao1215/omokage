# omokage demo tape (https://github.com/charmbracelet/vhs)
#
# Regenerate the GIF:
#   go build -o /tmp/omokage .          # build the CLI as `omokage`
#   rm -rf /tmp/omokage-demo && mkdir -p /tmp/omokage-demo
#   cp -r examples /tmp/omokage-demo/examples
#   PATH=/tmp:$PATH vhs doc/img/demo.tape   # writes demo.gif, then move it:
#   mv demo.gif doc/img/demo.gif
#
# vhs v0.11 parses absolute Output paths oddly, so the tape writes a relative
# demo.gif into the current directory and the GIF is moved into doc/img/ after.

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
