import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import ReactAce from "react-ace/lib/ace";
import { Badge, Button, Card, Col, Container, Form, ListGroup, Row } from "react-bootstrap";
import { Link, useNavigate, useParams } from "react-router-dom";

const ProviderConfigView = () => {
  const navigate = useNavigate()
  const { name } = useParams()
  const [providerConfig, setProviderConfig] = useState<any>()
  const [type, setType] = useState('')
  const [source, setSource ] = useState('')
  const [config, setConfig ] = useState('{}')
  const [secretRef, setSecretRef] = useState('')


  const getItem = useCallback(async () => {
    const res = await fetch(`http://localhost:8080/providers-configs/${name}`)
    const item = await res.json()

    setProviderConfig(item)
  }, [])

  useEffect(() => {
    getItem()
  }, [])

  useEffect(() => {
    setType(providerConfig?.type)
    setSource(providerConfig?.source)
    setConfig(JSON.stringify(providerConfig?.awsConfig, null, 2))
    setSecretRef(`${providerConfig?.secretRef?.namespace}/${providerConfig?.secretRef?.name}`)
  }, [providerConfig])
  
  return (
    <div style={{padding: '80px'}}>
      <h1>{providerConfig?.name}</h1>
      <Form.Group className="mb-3">
        <Form.Label>Type</Form.Label>
        <Form.Select disabled value={type} onChange={e => setType(e.target.value)}>
          <option value="" disabled>Select a type</option>
          {['AWS', 'Azure'].map(i => (
            <option value={i}>{i}</option>
          ))}
        </Form.Select>
        <small>Provider type</small>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Source</Form.Label>
        <Form.Select disabled value={source} onChange={e => setSource(e.target.value)}>
          <option value="" disabled>Select a type</option>
          {['SECRET', 'CURRENT_ACCOUNT'].map(i => (
            <option value={i}>{i}</option>
          ))}
        </Form.Select>
        <small>Credentials source</small>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Config</Form.Label>
        <ReactAce
          mode="json"
          theme="monokai"
          width='100%'
          height='200px'
          value={config}
          onChange={e => setConfig(e)}
          readOnly={true}
          enableBasicAutocompletion={true}
        />
        <small>Provider config</small>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Secret ref</Form.Label>
        <Form.Control type="text" disabled placeholder="type the plugin ref" value={secretRef} onChange={e => setSecretRef(e.target.value)} />
        <small>Secret ref </small>
      </Form.Group>
    </div>
  )
}

export default ProviderConfigView