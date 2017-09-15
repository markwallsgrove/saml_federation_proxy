import EventEmitter from 'eventemitter3';
import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';
import {
    Button, Modal, FormGroup, InputGroup,
    DropdownButton, MenuItem, FormControl,
    Accordion, Panel, Glyphicon
} from 'react-bootstrap';

import 'bootstrap/dist/css/bootstrap.css';
import './css/style.css';

const EE = new EventEmitter();

function paginateEntityDescriptors() {
    const url = '/1/api/entitydescriptors';

    return axios.get(url).then(resp => {
        // TODO: error?
        return resp.data;
    });
}

function entityDescriptor(id) {
    const url = '/1/api/entitydescriptor?id=' + id;

    return axios.get(url).then(resp => {
        // TODO: error?
        return resp.data;
    });
}

function createExport(name) {
    const url = `/1/api/exports`;
    const body = {
        name,
        entitydescriptor: []
    };

    // TODO: error
    return axios.post(url, body);
}

function getExports() {
    const url = `1/api/exports`;

    return axios.get(url).then(resp => {
        // TODO: error
        return resp.data;
    });
}

function activateEdOnExport(exportName, entityId) {
    const url = `1/api/exports/${exportName}`;
    const body = {
        "EntityDescriptors": {
            "Append": [entityId],
        }
    };

    return axios.patch(url, body).then(resp => {
        // TODO: error
        return resp.data;
    });
}

function deactiveEdOnExport(exportName, entityId) {
    const url = `1/api/exports/${exportName}`;
    const body = {
        "EntityDescriptors": {
            "Delete": [entityId],
        }
    };

    return axios.patch(url, body).then(resp => {
        // TODO: error
        return resp.data;
    });
}

class ExportEntityDescriptors extends React.Component {
    constructor(props) {
        super(props);

        // TODO: should a msg to create a export if none exist
        this.state = {
            entitydescriptors: [],
            exp: {
                EntityDescriptors: [],
                Name: '',
            },
        };

        EE.on('show-export-event', this.loadExportEvent, this);
    }

    loadExportEvent(exp) {
        console.log('>>>>>>>>>>>>>>>', exp);
        this.setState({exp});
    }

    toggleEdOnExport(entityId, isActive) {
        const exportName = this.state.exp.Name;
        const action = isActive
            ? deactiveEdOnExport(exportName, entityId)
            : activateEdOnExport(exportName, entityId);

        action.then(exp => {
            EE.emit('updated-export-event', exp);
            this.setState({exp});
        });
    }

    componentDidMount() {
        paginateEntityDescriptors().then(resp => {
            this.setState({entitydescriptors: resp});
        });
    }

    render() {
        const selected = this.state.exp.EntityDescriptors;
        const exportName = this.state.exp.Name;

        return (
            <div className="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main" id="entityDescriptorList">
                <Accordion>
                    {this.state.entitydescriptors.map((ed, i) => {
                        const isActive = selected.indexOf(ed.EntityID) > -1
                        const statusClass = isActive ? 'active' : 'inactive';
                        const statusIcon = isActive
                            ? 'glyphicon glyphicon-remove-circle'
                            : 'glyphicon glyphicon-ok-circle';
                        const btnClass = isActive ? 'btn-danger' : 'btn-success';
                        const btnText = isActive ? 'remove' : 'add';
                        const toggleEntityId = this.toggleEdOnExport.bind(this, ed.EntityID, isActive);

                        return <Panel header={ed.EntityID} eventKey={i} className={statusClass}>
                            EntityID: {ed.EntityID}
                            <br />
                            Federtion: {ed.FederationID}
                            <br /><br />
                            <Button onClick={toggleEntityId} className={btnClass}><Glyphicon glyph={statusIcon} /> {btnText}</Button>
                        </Panel>
                    })}
                </Accordion>
            </div>
        );
    }
}

class CreateExportInput extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showCreateExport: false,
        };
    }

    btnCreateExport() {
        createExport(this.input.value).then((resp) => {
            this.input.value = "";
            EE.emit('refresh-exports');
        });
    }

    render() {
        const createExport = this.btnCreateExport.bind(this);

        return (
            <div>
                <FormGroup>
                  <InputGroup>
                    <FormControl type="text" inputRef={(ref) => {this.input = ref}} />
                    <InputGroup.Button>
                      <Button onClick={createExport}>create</Button>
                    </InputGroup.Button>
                  </InputGroup>
                </FormGroup>
            </div>
        );
    }
}

class ExportSideBar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            exports: {},
            selected: '',
        };

        EE.on('refresh-exports', this.refreshExports, this);
        EE.on('updated-export-event', this.updateExportEvent, this);
    }

    componentDidMount() {
        this.refreshExports(true);
    }

    refreshExports(triggerEvent) {
        getExports().then(exports => {
            const mappedExports = {};
            exports.forEach(exp => {
                mappedExports[exp.Name] = exp;
            });

            this.setState({exports: mappedExports});

            if (triggerEvent && exports.length > 0) {
                this.triggerShowExportEvent(exports[0]);
            }
        });
    }

    triggerShowExportEvent(exp) {
        this.setState({selected: exp.Name});
        EE.emit('show-export-event', exp);
    }

    updateExportEvent(exp) {
        const exports = this.state.exports;
        exports[exp.Name] = exp;
        this.setState({exports});
    }

    render() {
        const trigger = this.triggerShowExportEvent;
        const exports = this.state.exports;

        return (
            <div id="exportSideBar" className="col-sm-3 col-md-2 sidebar">
                <CreateExportInput />
                <ul className="nav nav-sidebar">
                    {Object.values(exports).map(exp => {
                        const selected = exp.Name == this.state.selected ? 'active' : '';
                        return <li class="nav-link {selected}" onClick={trigger.bind(this, exp)}>{exp.Name}</li>
                    })}
                </ul>
            </div>
        );
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div className="row">
                <ExportSideBar />
                <ExportEntityDescriptors />
            </div>
        );
    }
}

ReactDOM.render(<App />,  document.getElementById('app'))
