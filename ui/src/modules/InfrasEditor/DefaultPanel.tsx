import React, { useCallback, useEffect, useState } from 'react'
import './index.css'
import { Alert, Button, ButtonGroup, Card, Form, ListGroup, ListGroupItem, Nav, ToggleButton } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { useSearchParams } from 'react-router-dom'

const getClassNameByExecution = (execution: any) => {
  if (execution?.status?.status === "SUCCESS") {
    return 'shared-infra-diagram__default-panel__execution--success'
  }

  if (execution?.status?.status === "RUNNING") {
    return 'shared-infra-diagram__default-panel__execution--running'
  }

  return 'shared-infra-diagram__default-panel__execution'
}

const DefaultPanel = ({ infra, onSave, onChange, goToView }: any) => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [name, setName] = useState(infra?.name || '')
  const [description, setDescription] = useState(infra?.description || '')
  const [providerConfigRef, setProviderConfigRef] = useState<any>(infra?.providerConfigRef || '')

  const [providersConfig, setProvidersConfig] = useState<any>([])

  const getProvidersConfigs = useCallback(async () => {
    const infraRes = await fetch(`http://localhost:8080/providers-configs`)
    const infra = await infraRes.json()

    setProvidersConfig(infra.items)
  }, [])

  useEffect(() => {
    getProvidersConfigs()
  }, [])

  useEffect(() => {
    setName(infra?.name)
    setDescription(infra?.description)
    setProviderConfigRef(infra?.providerConfigRef)
  }, [infra])

  useEffect(() => {
    onChange({
      name: name || '',
      namespace: "default",
      description: description || '',
      providerConfigRef: providerConfigRef || {},
    })
  }, [name, description, providerConfigRef])
  
  const handleCreate = () => {
    console.log(providerConfigRef)
    onSave({
      name,
      namespace: "default",
      description,
      providerConfigRef,
    })
  } 

  return (
    <div className='shared-infra-diagram__default-panel'>
      <div>
        <div className="d-grid gap-2 mb-2">

       
        <ButtonGroup>
          <ToggleButton
            value={searchParams.get("view") || ''}
            onClick={e => setSearchParams({ view: 'CODE' })}
            checked={searchParams.get("view") === 'CODE'}
            variant='secondary'
            type='radio'
          >
            <FontAwesomeIcon style={{cursor: 'pointer'}} icon="code" />
          </ToggleButton>
          <ToggleButton
            value={searchParams.get("view") || ''}
            onClick={e => setSearchParams({ view: 'DIAGRAM' })}
            checked={searchParams.get("view") === 'DIAGRAM'}
            variant='secondary'
            type='radio'
          >
            <FontAwesomeIcon style={{cursor: 'pointer'}} icon="diagram-project"/>
          </ToggleButton>
        </ButtonGroup>
        </div>
        {infra && <FontAwesomeIcon onClick={goToView} className='mb-2' style={{cursor: 'pointer'}} icon="arrow-left" />}
        <Card.Title>{infra?.name}</Card.Title>
        <Form.Group className='mb-3'>
          <Form.Label>Name</Form.Label>
          <Form.Control type="text" placeholder='Type the of yout shared infra' value={name} onChange={e => setName(e.target.value)} />
        </Form.Group>
        <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1" >
          <Form.Label>Description</Form.Label>
          <Form.Control as="textarea" rows={3} value={description} onChange={e => setDescription(e.target.value) } />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Providers configs</Form.Label>
          <Form.Select value={providerConfigRef} onChange={e => setProviderConfigRef(e.target.value)}>
            <option value="" disabled>Select a provider config</option>
            {providersConfig.map((i: any) => (
              <option value={i}>{i.name}</option>
            ))}
          </Form.Select>
        </Form.Group>
        <div className="d-grid gap-2">
          <Button onClick={handleCreate}>Save shared infra</Button>
        </div>
      </div>
    </div>
  )
  
}

export default DefaultPanel