import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React from "react"
import { Breadcrumb, Button, ListGroup } from "react-bootstrap"
import { useLocation, useNavigate, useParams } from "react-router-dom"

const infras = [
  { id: '1', name: 'Task Output 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec placerat, enim a cursus tempus.' },
  { id: '2', name: 'Task Output 2', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec placerat, enim a cursus tempus.' },
]

const TasksOutputs = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { workspaceId } = useParams()

  return (
    <div className="m-4">
      <div>
        <Breadcrumb>
          {location.pathname.split('/').map((path, index) => {
            if (path === '') {
              return null
            }

            return (
              <Breadcrumb.Item key={index} href="#">
                {path}
              </Breadcrumb.Item>
            )
          })}
        </Breadcrumb>
      </div>
      <div className="mb-3 d-flex align-items-center justify-content-between">
        <h1>Tasks Outputs</h1>
        <Button><FontAwesomeIcon icon="plus" /> Create</Button>
      </div>
      <ListGroup variant="flush">
        {infras.map((infra, index) => (
          <ListGroup.Item action onClick={() => navigate(`/workspaces/${workspaceId}/infras/${infra.id}`)}>
            <div className="ms-2 me-auto">
              <div className="fw-bold">{infra?.name}</div>
              {infra?.description}
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  )
}

export default TasksOutputs