import React from 'react'
import './index.css'
import { Alert, Card, ListGroup, ListGroupItem, Nav } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const getClassNameByExecution = (execution: any) => {
  if (execution?.status?.status === "SUCCESS") {
    return 'shared-infra-diagram__default-panel__execution--success'
  }

  if (execution?.status?.status === "RUNNING") {
    return 'shared-infra-diagram__default-panel__execution--running'
  }

  return 'shared-infra-diagram__default-panel__execution'
}

const DefaultPanel = ({ sharedInfra, onSelectExecution, onClose }: any) => {
  return (
    <Card className='shared-infra-diagram__default-panel'>
      <div>
        <Card.Title>{sharedInfra?.name}</Card.Title>
        <p>{sharedInfra?.description}</p>
        <div>
          <strong>Executions</strong>
          {sharedInfra?.status?.executions?.map((i: any) => (
            <Card 
              className={`${getClassNameByExecution(i)} mt-1`}
              style={{cursor: 'pointer'}} 
              onClick={() => onSelectExecution(i)}
            >
              {i.name}
            </Card>
          ))}
          
        </div>
      </div>
    </Card>
  )
  
}

export default DefaultPanel