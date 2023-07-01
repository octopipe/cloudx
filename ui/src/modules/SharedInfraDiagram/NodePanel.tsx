import React, { useEffect, useState } from 'react'
import './index.css'
import { Alert, Card, Form, ListGroup } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import ReactAce from 'react-ace/lib/ace'

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";

const nodeDefaultTypes = [
  { value: 'terraform', label: 'Terraform' }
]

const nodeAwsTypes = [
  { value: 'eks', label: 'EKS' },
  { value: 'secrets-manager', label: 'Secrets Manager' },
  { value: 'ecr', label: 'ECR' },
  { value: 's3', label: 'S3' },
]

const getTypesByCategory = (category: any) => {
  return category === "default" ? nodeDefaultTypes : nodeAwsTypes
}

const NodePanel = ({ selectedNode, action, onClose, onChange }: any) => {
  const [currentNode, setCurrentNode] = useState()
  const [name, setName] = useState(selectedNode?.data?.name)
  const [type, setType] = useState(selectedNode?.data?.type)
  const [ref, setRef] = useState(selectedNode?.data?.ref)
  const [inputs, setInputs] = useState(JSON.stringify(selectedNode?.data?.inputs, null, 2))

  useEffect(() => {
    onChange({
      ...selectedNode,
      data: {
        ...selectedNode.data,
        label: name,
        name,
        type,
        ref,
        inputs: JSON.parse(inputs),
      }
    })
  }, [name, type, ref, inputs])

  return (
    <div className='shared-infra-diagram__node-panel'>
      <Card.Body>
        <div>
          <div className='d-flex justify-content-end'>
            <FontAwesomeIcon
              icon="close"
              style={{cursor: 'pointer'}}
              onClick={onClose}
            />
          </div>
          {(action === "EDIT" || action === "CREATE" ) && (
            <>
              <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
                <Form.Label>Name</Form.Label>
                <Form.Control type="text" placeholder="type the plugin name" value={name} onChange={e => setName(e.target.value)}/>
              </Form.Group>
              <Form.Group className="mb-3">
                <Form.Label>Type</Form.Label>
                <Form.Select defaultValue={type} onChange={e => setName(e.target.value)}>
                  <option value="" disabled>Select a type</option>
                  {getTypesByCategory(selectedNode?.data?.category).map(i => (
                    <option value={i.value}>{i.value}</option>
                  ))}
                </Form.Select>
              </Form.Group>
              {selectedNode?.data?.category !== "aws" && (
                <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
                  <Form.Label>Ref</Form.Label>
                  <Form.Control type="text" placeholder="type the plugin ref" value={ref} onChange={e => setRef(e.target.value)} />
                </Form.Group>
              )}
              <Form.Group className="mb-3">
                <Form.Label>Inputs</Form.Label>
                <ReactAce
                  mode="json"
                  theme="monokai"
                  width='100%'
                  height='200px'
                  value={inputs}
                  onChange={e => setInputs(e)}
                  enableBasicAutocompletion={true}
                />
              </Form.Group>
            </>
          )}

          {(action !== "EDIT" && action !== "CREATE" ) && (
            <>
              {selectedNode?.data?.error && <Alert variant='danger'>{selectedNode?.data?.error}</Alert>}
              <Card.Title>{selectedNode?.data?.name}</Card.Title>
              <div><strong>Ref: </strong>{selectedNode?.data?.ref}</div>
              <div><strong>Type: </strong>{selectedNode?.data?.type}</div>
              <div>
                <strong>Plugins:</strong><br/>
                <ListGroup className='mt-2'>
                  {selectedNode?.data?.inputs?.map((i: any) => (
                    <ListGroup.Item style={{border: '1px solid #ccc', padding: '5px'}}>
                      <strong>{i?.key}: </strong>{i?.value}
                    </ListGroup.Item>
                  ))}
                </ListGroup>
              </div>
            </>
          )}
        </div>
      </Card.Body>
    </div>
  )
  
}

export default NodePanel