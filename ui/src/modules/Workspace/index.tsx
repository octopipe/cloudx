import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import React, { useEffect } from "react"
import { Breadcrumb, Dropdown, Nav } from "react-bootstrap"
import { NavLink, Outlet, useLocation, useParams } from "react-router-dom"

const menuItems = [
  { icon: 'diagram-project', to: 'infras' },
  { icon: 'circle-nodes', to: 'tasks-outputs' },
  { icon: 'cloud', to: 'providers-config' },
  { icon: 'tower-broadcast', to: 'webhooks' },
  { icon: 'box', to: 'remote' },

]

const Workspace = () => {
  const location = useLocation()
  const { workspaceId } = useParams()

  return (
    <div style={{ height: '100vh', display: 'flex' }}>
      <div className="d-flex flex-column justify-content-between text-white bg-dark" style={{width: '4.5rem'}}>
        <Nav variant="pills" activeKey={location.pathname} className="nav-flush flex-column mb-auto text-center">
          <div style={{height: '60px', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
            CX
          </div>
          {menuItems.map((item, index) => (
            <NavLink to={`/workspaces/${workspaceId}/${item?.to}`} className="nav-link text-white py-3">
              <FontAwesomeIcon size="lg" icon={item?.icon as any} />
            </NavLink>
          ))}
        </Nav>
        <Nav variant="pills" className="nav-flush flex-column text-center">
          <NavLink to={`/workspaces/${workspaceId}/settings`} className="nav-link text-white py-3">
            <FontAwesomeIcon size="lg" icon="gear" />
          </NavLink>
        </Nav>
      </div>
      <div style={{width: '100%'}}>
        <Outlet />
      </div>
    </div>
  )
}

export default Workspace