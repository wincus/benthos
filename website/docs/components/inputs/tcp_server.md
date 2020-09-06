---
title: tcp_server
type: input
deprecated: true
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/input/tcp_server.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

DEPRECATED: This component is deprecated and will be removed in the next major
version release. Please consider moving onto [alternative components](#alternatives).

```yaml
# Config fields, showing default values
input:
  tcp_server:
    address: 127.0.0.1:0
    delimiter: ""
    max_buffer: 1e+06
    multipart: false
```

Creates a server that receives messages over TCP. Each connection is parsed as a
continuous stream of line delimited messages.

If multipart is set to false each line of data is read as a separate message. If
multipart is set to true each line is read as a message part, and an empty line
indicates the end of a message.

If the delimiter field is left empty then line feed (\n) is used.

The field `max_buffer` specifies the maximum amount of memory to
allocate _per connection_ for buffering lines of data. If a line of data from a
connection exceeds this value then the connection will be closed.

