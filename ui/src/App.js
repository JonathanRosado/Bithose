import React, {useEffect, useRef, useState} from 'react';
import logo from './logo.svg';
import './App.css';

const subscribeText =
`{ 
  "type": "subscribe",
  "criteria": [
    { "operator": "==", "label_pair": { "name": "channel", "value": "cool_channel" } },
    { "operator": ">", "label_pair": { "name": "num_of_chars", "value": 5 } }
  ] 
}
`
const messageText =
`{
  "type": "message",
  "message": {
    "body": "hello",
    "label_pairs": [
      { "name": "channel", "value": "cool_channel" },
      { "name": "num_of_chars", "value": 12 }
    ]
  }
}
`

function TextWidget({t}) {
    const ws = useRef(null)
    const messagesEnd = useRef(null)
    const [messages, setMessages] = useState([])
    const [textarea, setTextarea] = useState("")

    const onSend = () => {
        console.log(`sending\n${textarea}`)
        console.log(JSON.parse(textarea))
        ws.current.send(textarea)
    }

    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:80")
        ws.current.onclose = () => setMessages([...messages, "websocket closed"])
        ws.current.onopen = () => setMessages([...messages, "websocket open"])
        ws.current.onerror = err => setMessages([...messages, "there's been an error: " + err.toString()])
        return () => ws.current.close()
    }, [])

    useEffect(() => {
        if (!ws.current) return;
        ws.current.onmessage = msg => setMessages([...messages, msg.data])
    })

    useEffect(() => {
        setTimeout(() => {
            messagesEnd.current.scrollIntoView({ behavior: "smooth" })
        }, 200)
    })

    return (
        <div className={`t${t} textbox`}>
            <div className="text">
                {messages.map(m => <div style={{fontFamily: "\"Courier New\", Courier, monospace", fontSize: 14}}>{m}</div>)}
                <div style={{ float:"left", clear: "both" }}
                     ref={messagesEnd}>
                </div>
            </div>
            <div className="send">
                <textarea id={'textarea'} className="textarea"
                          style={{fontFamily: "\"Courier New\", Courier, monospace", fontSize: 11}}
                          onChange={e => setTextarea(e.target.value)} value={textarea}/>
                <div style={{display: "flex", flexDirection: "column"}}>
                    <button onClick={() => setTextarea(subscribeText)}>subscribe</button>
                    <button onClick={() => setTextarea(messageText)}>message</button>
                    <button onClick={onSend} style={{flex: 1}}>send</button>
                </div>
            </div>
        </div>
    )
}

function App() {
  return (
    <div className="App">
        {([1,2,3,4,5,6].map(t => <TextWidget t={t} /> ))}
    </div>
  );
}

export default App;
