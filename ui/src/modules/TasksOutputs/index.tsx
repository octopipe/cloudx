import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React, { useCallback, useEffect, useState } from "react"
import { Breadcrumb, Button, ListGroup } from "react-bootstrap"
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { useFetch } from "use-http"

const taskOutputs = [
  { id: '1', name: 'Task Output 1', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec placerat, enim a cursus tempus.' },
  { id: '2', name: 'Task Output 2', description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec placerat, enim a cursus tempus.' },
]

const TasksOutputs = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { workspaceId } = useParams()
  const { patch, get, response, loading, error } = useFetch({ cachePolicy: 'no-cache' as any })
  const [taskOutputs, setTaskOutputs] = useState<any>([])

  useEffect(() => {
    getTaskOutputs()
    const interval = setInterval(() => {
      getTaskOutputs()
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const getProgressNow = (taskOutput: any) => {
    return  (taskOutput?.status?.tasks?.length * 100) / taskOutput?.tasks?.length
  }
  
  const getTaskOutputs = useCallback(async () => {
    const taskOutputs = await get(`/connections-interfaces`)
    if (response.ok) setTaskOutputs(taskOutputs?.items)
  }, [get])

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
        {taskOutputs.map((taskOutput: any, index: any) => (
          <ListGroup.Item key={index}>
            <div className="ms-2 me-auto">
              <Breadcrumb>
                <Breadcrumb.Item active>
                  Infra
                </Breadcrumb.Item>
                <Breadcrumb.Item onClick={() => navigate(`/workspaces/${workspaceId}/infras/${taskOutput?.infra?.name}/last-execution`)}>
                  {taskOutput?.infra?.name}
                </Breadcrumb.Item>
                <Breadcrumb.Item active>
                  Task
                </Breadcrumb.Item>
                <Breadcrumb.Item onClick={() => navigate(`/workspaces/${workspaceId}/infras/${taskOutput?.infra?.name}/last-execution/task/${taskOutput?.taskName}`)}>
                  {taskOutput?.taskName}
                </Breadcrumb.Item>
                <Breadcrumb.Item active>
                  Task output
                </Breadcrumb.Item>
                <Breadcrumb.Item href="#">
                  {taskOutput?.name}
                </Breadcrumb.Item>
              </Breadcrumb>
              {taskOutput?.description}
            </div>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  )
}

export default TasksOutputs