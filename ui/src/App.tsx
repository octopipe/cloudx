import React from 'react'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Main from './modules/Main'
import Infras from './modules/Infras'
import ProvidersConfig from './modules/ProviderConfig'
import InfraView from './modules/InfrasView'
import TaskOutputs from './modules/TaskOutputs'
import './core/components/icon'
import InfraEditor from './modules/InfrasEditor'
import ProviderConfigView from './modules/ProviderConfigView'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />}>
          <Route path='' element={<Navigate to="/infra" replace />} />
          <Route path='/infra' element={<Infras />} />
          <Route path='/infra/create' element={<InfraEditor />}/>
          <Route path='/infra/:name/edit' element={<InfraEditor />}/>
          <Route path='/infra/:name' element={<InfraView />}/>
          <Route path='/connection-interfaces' element={<TaskOutputs />}/>
          <Route path='/providers-config' element={<ProvidersConfig />}/>
          <Route path='/providers-config/:name' element={<ProviderConfigView />}/>
        </Route>
        <Route path="login" element={<div>Login</div>} />
        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
