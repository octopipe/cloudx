import React from 'react'
import './index.css'
import { Card, ListGroup } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const NodePanel = ({ node, onClose }: any) => {
  return (
    <div className='shared-infra-diagram__node-panel'>
      <Card.Body>
        <div className='d-flex justify-content-end'>
          <FontAwesomeIcon
            icon="close"
            style={{cursor: 'pointer'}}
            onClick={onClose}
          />
        </div>
        <Card.Title>{node?.data?.name}</Card.Title>
        <div><strong>Ref: </strong>{node?.data?.ref}</div>
        <div><strong>Type: </strong>{node?.data?.type}</div>
        <div>
          <strong>Plugins:</strong><br/>
          <ListGroup className='mt-2'>
            {node?.data?.inputs?.map((i: any) => (
              <ListGroup.Item style={{border: '1px solid #ccc', padding: '5px'}}>
                <strong>{i?.key}: </strong>{i?.value}
              </ListGroup.Item>
            ))}
          </ListGroup>
        </div>
      </Card.Body>
    </div>
  )
  
}

export default NodePanel