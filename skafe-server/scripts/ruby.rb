#!/usr/bin/ruby
require 'json'

def mainLoop

	event = nil

	# Main run loop
	loop do
		# Get the command from stdin
		fullCmd = gets()

		return if fullCmd == nil

		# Split into parts
		commands = fullCmd.split(' ')

		case commands[0] 
			
			when /EVENT/i
				# Handle New Event
				if commands.length != 2 then
					puts "ERR invalid command format"
					next
				end
				event = JSON.parse(commands[1])

			when /CMD/i
				# Handle Command
				if event == nil then
					puts "ERR No event set"
					next
				end
				
				begin
					send(commands[1], event)
				rescue
					puts "ERR No such method"
				end
			else
				puts "ERR Command not understood"

		end

	end
end

def requireAll
	Dir[File.dirname(__FILE__) + '/lib/*.rb'].each do |file| 
  		require File.basename(file, File.extname(file))
	end
end

# Load all the user scripts
requireAll

# Run the main loop
mainLoop
