# Localization

The Goyave framework provides a convenient way to support multiple languages within your application. Out of the box, Goyave only provides the `en-US` language.

Language files are stored in the `resources/lang` directory.

```
.
└── resources
    └── lang
        └── en-US (language name)
            ├── fields.json (optional)
            ├── locale.json (optional)
            └── rules.json (optional)
```

Each language has its own directory and should be named with an [ISO 639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) language code. You can also append a variant to your languages: `en-US`, `en-UK`, `fr-FR`, `fr-CA`, ... **Case is important.**

Each language directory contains three files. Each file is **optional**.
- `fields.json`: field names translations and field-specific rule messages.
- `locale.json`: all other language lines.
- `rules.json`: validation rules messages.

All directories in the `resources/lang` directory are automatically loaded when the server starts.

Learn more about localization [here](https://system-glitch.github.io/goyave/guide/advanced/localization.html).