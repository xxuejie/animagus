$LOAD_PATH.unshift(File.expand_path(File.dirname(__FILE__)))

require "grpc"
require "generic_services_pb"

def main
  stub = Generic::GenericService::Stub.new("127.0.0.1:4000", :this_channel_is_insecure)
  request = Generic::GenericParams.new(
    name: "balance",
    params: [
      Ast::Value.new(
        t: Ast::Value::Type::BYTES,
        raw: [ARGV[0] || ""].pack("H*")
      )
    ]
  )
  response = stub.call(request)
  p response
end

main
