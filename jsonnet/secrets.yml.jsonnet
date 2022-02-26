{
  apiVersion: 'v1',
  stringData: {
    DB_CONNECTION_STRING: std.extVar('DB_CONNECTION_STRING'),
    BUCKET_ENDPOINT: std.extVar('BUCKET_ENDPOINT'),
    BUCKET_ACCESS_KEY: std.extVar('BUCKET_ACCESS_KEY'),
    BUCKET_SECRET_KEY: std.extVar('BUCKET_SECRET_KEY'),
    BUCKET_NAME: std.extVar('BUCKET_NAME'),
    BUCKET_LOCATION: std.extVar('BUCKET_LOCATION'),
  },
  kind: 'Secret',
  metadata: {
    name: 'oceanus-secrets',
    namespace: 'oceanus',
  },
}
