#!/usr/bin/env ruby

ARGF.each do |line|
  p line.gsub(/\{\s"\$date" : "([0-9\-:T\+\.]+)" \}/, '"\1"')
end
