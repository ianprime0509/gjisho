#!/usr/bin/env python3
# count-strokes.py - count the number of strokes in a list of input characters

# This script requires the Unihan_IRGSources.txt file available in Unihan.zip,
# which can be downloaded from https://www.unicode.org/Public/UCD/latest/ucd/.

import argparse
import re
import sys

parser = argparse.ArgumentParser()
parser.add_argument("unihan_irg")
args = parser.parse_args()

stroke_re = re.compile(r"U\+([0-9A-F]+)\tkTotalStrokes\t([0-9]+)")
# There are several radicals and other components not in the Unihan database,
# and their stroke orders must be provided manually
stroke_counts = {
    "⺅": 2,
    "⺉": 2,
    "ユ": 2,
    "⺌": 3,
    "⺖": 3,
    "⺡": 3,
    "⺨": 3,
    "⺾": 3,
    "⻌": 3,
    "⻏": 3,
    "⺣": 4,
    "⺹": 4,
    "⽱": 5,
    "⽧": 5,
    "⺲": 5,
    "⻂": 5,
}
with open(args.unihan_irg) as irg:
    for line in irg:
        match = stroke_re.match(line)
        if match is not None:
            c = chr(int(match.group(1), 16))
            strokes = int(match.group(2))
            stroke_counts[c] = strokes

for c in sys.stdin:
    c = c.strip()
    print("{}\t{}".format(c, stroke_counts.get(c, -1)))
