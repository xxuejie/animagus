$LOAD_PATH.unshift(File.expand_path(File.join(File.dirname(__FILE__), "..", "..", "ruby")))

require "grpc"
require "generic_services_pb"

def bin_to_hex(bin)
  "0x#{bin.unpack1('H*')}"
end

def main
  stub = Generic::GenericService::Stub.new("127.0.0.1:4000", :this_channel_is_insecure)
  request = Generic::GenericParams.new(
    name: "nervosdao_deposits"
  )
  response = stub.stream(request)
  response.each do |r|
    puts "New NervosDAO deposit at tx hash: #{bin_to_hex(r.children[0].raw)}, index: #{r.children[1].u}"
  end
end

main
