import React from "react"
import { Button, Card, Col, Container, Row } from "react-bootstrap"
import { useNavigate } from "react-router-dom"

const workspaces = [
  { name: "Workspace 1" },
  { name: "Workspace 2" },
  { name: "Workspace 3" },
  { name: "Workspace 4" },
  { name: "Workspace 5" },
]

const Workspaces = () => {
  const navigate = useNavigate()

  return (
    <>
      <Container>
        <Row className="mt-4">
          {workspaces.map((workspace, idx) => (
            <Col xs={3} key={idx}>
              <Card className="mt-4">
                <Card.Header>{workspace?.name}</Card.Header>
                <Card.Body>
                  <div className="d-grid gap-2">
                    <Button 
                      variant="outline-primary"
                      onClick={() => navigate(`/workspaces/${idx + 1}`)}
                    >
                      Enter
                    </Button>
                  </div>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      </Container>
    </>
  )
}

export default Workspaces