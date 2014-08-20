# mongoexport_date_fix

It's a simple tool, to fix mongoexport weird json output, where date fields are stored like:

`{"date_field" : { "$date" : "2013-02-19T15:15:20.079+0000"}}`

instead:

`{"date_field" : "2013-02-19T15:15:20.079+0000"}`

Usage:

`mongoexport --collection collection_name | mongoexport_date_fix | gzip -9 > collection_name.json.gz`
