source = ["./bin/macos/universe"]
bundle_id = "io.heeser.meta-cli"

apple_id {
  username = "yann@heeser.io"
  password = "@env:AC_PASSWORD"
  provider = "S5H7RTPRN5"
}

sign {
  application_identity = "C9C0D18A16BEE8034D9601F247DFD2A5D24926DB"
}

zip {
  output_path = "bin/macos-x64-signed.zip"
}
