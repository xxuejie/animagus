$LOAD_PATH.unshift(File.expand_path(File.join(File.dirname(__FILE__), "..", "..", "ruby")))

require "json"
require "grpc"
require "generic_services_pb"

if ARGV.length != 4
  puts "Usage: ruby transfer.rb <udt type arg> <from lock arg> <to lock arg> <amount>"
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
    name: "transfer",
    params: [
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: hex_to_bin(ARGV[0])
      ),
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: hex_to_bin(ARGV[1])
      ),
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: hex_to_bin(ARGV[2])
      ),
      Ast::Value.new(
        t: Ast::Value::Type::UINT64,
        u: ARGV[3].to_i
      ),
    ]
  )
  response = stub.call(request)
  p response
  puts JSON.pretty_generate(JSON.parse(response.raw))
end

main
