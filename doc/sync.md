# Funkapparat

## Sync

### Configuration

Besides the exporters in the [exporter configuration](configuration.md#exporter), you can configure Funkapparat to sync all changes in a [Redis stream](https://redis.io/docs/latest/develop/data-types/streams/).

Simply add `app.redis_url` to your config, e.g., `redis://localhost:3351/0`.

### Protocol

The events (sync, create, update, delete) will be according to the news-schema defined in the [Protokolapparat repository](https://github.com/potibm/protokolapparat).

