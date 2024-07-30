[dev setup documentation]: docs/dev-setup.md#requirements

# Contributing

Contributing encompasses repository specific requirements.

## Requirements

To review the requirements, please refer to the [dev setup documentation].

## Writing code

### Running locally

The recommended flow for development would be have two sessions, one for the "frontend":

```bash
npm run start:dev
```

And one for the "backend":

```bash
cd ../bff
docker compose -f docker-compose.yaml up
```

Once you have both services ready, you can open the dashboard locally at: `http://localhost:4010`. The dev server will reload automatically when you make changes.


## Debugging and Testing

[TBD]

## Environment variables

[TBD]

### Building your image

[TBD]