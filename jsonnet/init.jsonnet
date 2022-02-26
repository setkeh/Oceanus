local Secrets = import 'secrets.yml.jsonnet';
std.manifestYamlDoc(Secrets)
