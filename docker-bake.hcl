group "default" {
  targets = ["funkapparat"]
}

target "funkapparat" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = ["linux/amd64", "linux/arm64"]

  labels = {
    "org.opencontainers.image.url" = "https://github.com/potibm/funkapparat"
    "org.opencontainers.image.source" = "https://github.com/potibm/funkapparat"
    "org.opencontainers.image.documentation" = "https://github.com/potibm/funkapparat/tree/main/doc"
    "org.opencontainers.image.authors" = "potibm"
  }
  
  annotations = [
    "index,manifest:org.opencontainers.image.title=Funkapparat",
    "index,manifest:org.opencontainers.image.description=An editor for news and announcements at demoparties.",
    "index,manifest:org.opencontainers.image.url=https://github.com/potibm/funkapparat",
    "index,manifest:org.opencontainers.image.source=https://github.com/potibm/funkapparat",
    "index,manifest:org.opencontainers.image.documentation=https://github.com/potibm/funkapparat/tree/main/doc",
    "index,manifest:org.opencontainers.image.licenses=MIT",
    "index,manifest:org.opencontainers.image.authors=potibm"
  ]
}
