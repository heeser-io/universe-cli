source = ["./bin/macos-silicon/universe"]
bundle_id = "io.universecloud.cli"

apple_id {
  username = "yann@heeser.io"
  password = "@env:AC_PASSWORD"
  provider = "S5H7RTPRN5"
}

sign {
  application_identity = "9ED3857AFEC2BAAEED941122856E838450887B4B"
}

zip {
  output_path = "bin/macos-silicon-signed.zip"
}
