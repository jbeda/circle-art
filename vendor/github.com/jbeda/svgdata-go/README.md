# svgdata-go

**WARNING:** This is very early library that I wrote for myself and anything may change at any time.

A go library for representing SVG data and then writing it out.
This isn't the best API but it is useful.

The one thing this does well is collect a set of path segment and then construct a Path from those by looking for the segments that match up.
This is useful for CAD applications where often times you'll get a collection of path segments and need to assemble continuous paths from those segments.