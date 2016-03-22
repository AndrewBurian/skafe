#!/usr/bin/ruby
require 'json'

def mainLoop

	event = nil

	# Main run loop
	loop do
		# Get the command from stdin
		fullCmd = STDIN.gets

		return if fullCmd == nil

		# Split into parts
		commands = fullCmd.chomp.split(' ')

		case commands[0] 

			when /EVENT/i
				# Handle New Event
				if commands.length != 2 then
					STDOUT.puts "ERR invalid command format"
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
					result = send(commands[1], event)
					if result then
						STDOUT.puts "RESP True"
					else
						STDOUT.puts "RESP False"
					end
				rescue
					STDOUT.puts "ERR No such method"
				end
			else
				STDOUT.puts "ERR Command not understood"

		end

	end
end

def requireAll(libPath)

	Dir[libPath + '*.rb'].each do |file|
		require libPath + File.basename(file, File.extname(file))
	end
end

# Load all the user scripts
requireAll ARGV[0]

# Run the main loop
mainLoop
