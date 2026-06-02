# omokage --explain demo tape (https://github.com/charmbracelet/vhs)
#
# Regenerate with `make demo`, which builds the CLI, seeds a throwaway project,
# runs this tape, and moves explain.gif into doc/img/. The init and train steps
# are hidden so the GIF focuses on the late-stage tuning view: what to fix, and
# where. vhs v0.11 parses absolute Output paths oddly, so the tape writes a
# relative explain.gif and `make demo` moves it into doc/img/.

Output explain.gif

Set Shell bash
Set FontSize 16
Set Width 1200
Set Height 1040
Set Padding 24
Set Theme "Catppuccin Mocha"

Hide
Type "export PATH=/tmp:$PATH; cd /tmp/omokage-demo" Enter
Type "omokage init >/dev/null; omokage train --author me examples/ja/posts >/dev/null; clear" Enter
Show

Type "# Late-stage tuning on a Japanese corpus: what to fix, and which paragraph?" Enter
Type "omokage check --author me --explain examples/ja/draft-keeps-voice.md" Enter
Sleep 4500ms
