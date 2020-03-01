Conventions to help cover [the first specs of DSE](https://docs.google.com/document/d/112-mwqZoNb4LwTcbWvWABWeHBDoTQNaxXvmqwGDbMUU/edit#)

## Conventions

- snake_case in DB for attribute names
- CamelCase (first letter uppercase) for collection/table names

Notes :
- the MessageVersions could be embedded in the Message, or available as a separate collection (and we need a foreign key ID)
- Some attributes could be denormalized directly in the DB for easier querying (or we could use some cache system, or not denormalize at all, TBD)
- This is only meant to be pseudo code, as long as you understand it's fine

## Collections

Message
- `id`
- `user_id` string - required (Slack user ID)
- `created_at` Date - required (date of creation in DB)
- `written_at` Date - required (date when first message version was written)
- `slack_message_id` string - required with index (slack message_id)
- `versions` []MessageVersion - required (there should always be at least one version)
- `all_messages_analyzed` boolean - required default false (true if all versions analyzed - denormalized attribute)
- `last_message_quality` - optional default nil (denormalized attribute))
- `improved` - boolean or nil default nil (at least one message_version improved regarding a previous one)
- `number_of_edits_before_improved` integer - optional (number of messages separating the first bad version from the last good version. Leave nil if there are no bad message versions)

MessageVersion
- `id`
- `created_at` - required (date of creation in DB)
- `written_at` - required (date when message version was written)
- `text` String - required
- `analyzed` Boolean
- `quality` Float - required when Sentiment is defined
- `sentiment` Map{Negative: Float, Positive:Float} optional
- `bad_quality_notified_on_slack` Boolean, default nil

## Services

(names are only given as an example, feel free to come up with something more explicit)

Service ProcessIncomingUserMessage
- If the Message (identified by messageId) does not exist in DB, creates the Message and set fields to default
- Update the message `all_messages_analyzed` and `lastMessageQuality`
- Creates the MessageVersion and set the fields to default

Service ProcessNLPResults
- Change a message_version `analyzed`, `sentiment`, and sets the resulting `quality`
- updates a message global stats `lastMessageQuality`, `all_messages_analyzed`, `improved`, `improved`, `number_of_edits_before_improved`
- if the quality is under the given threshold, trigger `NotifyUserOfBadMessageQuality`

Service NotifyUserOfBadMessageQuality
- Set `bad_quality_notified_on_slack` to true on the associated messageVersion after a successful POST

Service HomeStats
- # of Total des messages Analyzed par DSE --> `MessageVersion.count`
- # messages that have at least one message version with a bad Quality
  - if embedded : `Message.where("at least one 'versions.quality' lesser than X")`
  - if separate collection : `MessageVersion.where("quality lesser than => X").distinct("message_id").count`
- % of messages with bad quality
  - if embedded : `Message.where("every 'versions.quality' lesser than X")`
- % of improved messages : Message.where(improved: true)
- Number of iterations to get to a good quality message : `Message.avg("number_of_edits_before_improved")`

