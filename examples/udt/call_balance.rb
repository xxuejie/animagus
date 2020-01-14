$LOAD_PATH.unshift(File.expand_path(File.join(File.dirname(__FILE__), "..", "..", "ruby")))

require "grpc"
require "generic_services_pb"

if ARGV.length != 2
  puts "Usage: ruby call_balance.rb <udt type arg> <lock arg>"
  exit 1
end

def hex_to_bin(hex)
  hex = hex[2..-1] if hex.start_with?("0x")
  [hex].pack("H*")
end

def unpack_amount(data)
  values = data.unpack("Q<Q<")
  (values[1] << 64) | values[0]
end


def main
  stub = Generic::GenericService::Stub.new("127.0.0.1:4000", :this_channel_is_insecure)
  request = Generic::GenericParams.new(
    name: "balance",
    params: [
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: hex_to_bin(ARGV[0])
      ),
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: hex_to_bin(ARGV[1])
      ),
    ]
  )
  response = stub.call(request)
  p response
  puts "Amount: #{unpack_amount(response.raw)}"
end

main
