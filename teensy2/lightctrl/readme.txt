==== Light Ctrl ====

* Should appear as two USB serial devices /dev/ttyACM?

=== ttyACM0 ===

accepts input in the form:
'>\x12\x34\x56'
and thus sends the rf433 power outlet control sequence \x12\x34\x56 over the air

=== ttyACM1 ===

accepts input in the form of one char & 0x3F where each bit set corresponds
to the desired ceiling lights state. (1==on 0==off)

sends output in the form "%c%c%c\n" on any status change.

The first char is the current state in the same format as the input described above.

The following two chars encode the last pressed buttons:
Each set bit (MSB on the left aka first) signifies that one of the buttons has been pressed.

Starting from the upper left to the lower right on the button panel:
bit00: upper left button pressed UP / ON
bit01: upper left button pressed DOWN / OFF
bit02: upper middle button pressed UP / ON
bit03: upper middle button pressed DOWN / OFF
....
....
bit12: lower left single push button pressed
bit13: lower middle single push button pressed
bit14: lower right single push button pressed
