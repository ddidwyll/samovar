let store = ServerStateStore.make(
  ~url="//localhost:4000/store/feed",
  ~decodeValue=JSON.Decode.string,
)
