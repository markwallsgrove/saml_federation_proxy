import EventEmitter from 'eventemitter3';
import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';

import {
    Button, Modal, FormGroup, InputGroup,
    DropdownButton, MenuItem, FormControl,
    Accordion, Panel, Glyphicon
} from 'react-bootstrap';

import {
  BrowserRouter as Router,
  Route,
  Link,
  Switch
} from 'react-router-dom'

import 'bootstrap/dist/css/bootstrap.css';
import './css/style.css';

const EE = new EventEmitter();

// function paginateEntityDescriptors() {
//     const url = '/1/api/entitydescriptors';

//     return axios.get(url).then(resp => {
//         // TODO: error?
//         return resp.data;
//     });
// }

// function entityDescriptor(id) {
//     const url = '/1/api/entitydescriptor?id=' + id;

//     return axios.get(url).then(resp => {
//         // TODO: error?
//         return resp.data;
//     });
// }

// function createExport(name) {
//     const url = `/1/api/exports`;
//     const body = {
//         name,
//         entitydescriptor: []
//     };

//     // TODO: error
//     return axios.post(url, body);
// }

// function getExports() {
//     const url = `1/api/exports`;

//     return axios.get(url).then(resp => {
//         // TODO: error
//         return resp.data;
//     });
// }

// function activateEdOnExport(exportName, entityId) {
//     const url = `1/api/exports/${exportName}`;
//     const body = {
//         "EntityDescriptors": {
//             "Append": [entityId],
//         }
//     };

//     return axios.patch(url, body).then(resp => {
//         // TODO: error
//         return resp.data;
//     });
// }

// function deactiveEdOnExport(exportName, entityId) {
//     const url = `1/api/exports/${exportName}`;
//     const body = {
//         "EntityDescriptors": {
//             "Delete": [entityId],
//         }
//     };

//     return axios.patch(url, body).then(resp => {
//         // TODO: error
//         return resp.data;
//     });
// }

class Navbar extends React.Component {
    render() {
      return (
          <div id="navbar" className="collapse navbar-collapse">
              <ul className="nav navbar-nav">
                <li><Link to="/">Home</Link></li>
                <li><Link to="/idps">IDPs</Link></li>
                <li><Link to="/federations">Federations</Link></li>
                <li><Link to="/groups">Groups</Link></li>
                <li><Link to="/exports">Exports</Link></li>
              </ul>
          </div>
      );
    }
}

class Sidebar extends React.Component {
    constructor(props) {
        super(props);
    }

    onClick(link) {
        this.props.handleSelection(link);
    }

    render() {
        return (
            <div className="col-sm-2 sidebar">
                <ul>
                {this.props.links.map((link) => {
                    const stateClass = this.props.selected === link.id
                        ? 'active'
                        : 'inactive';

                    return <li
                        key={link.id}
                        className={stateClass}
                        onClick={this.onClick.bind(this, link)}>
                            {link.name}
                    </li>;
                })}
                </ul>
            </div>
        );
    }
}

class Home extends React.Component {
    render() {
        return (
            <div id="home">
                <Sidebar links={[]}/>
                <div className="col-md-10">home</div>
            </div>
        );
    }
}

class ServiceDetails extends React.Component {
    constructor(props) {
        super(props);
    }

    submit() {
        this.props.submit(this.service());
    }

    delete() {
        this.props.delete(this.service());
    }

    handleChange() {
        this.props.handleChange(this.service());
    }

    service() {
        return {
            name: this.refs.name.value,
            endpoint: this.refs.endpoint.value,
            certificateLocation: this.refs.certificateLocation.value,
            id: this.props.service.id,
        };
    }

    render() {
        const deleteButton = this.props.service.id && this.props.service.id !== '1'
            ? <button type="submit"
                className="btn btn-danger"
                onClick={this.delete.bind(this)}>delete</button>
            : null;

        return (
            <div className="col-sm-10 serviceDetails">
                <div className="form-group">
                    <label htmlFor="name">Name</label>
                    <input
                        className="form-control"
                        type="text"
                        placeholder="name"
                        name="name"
                        ref="name"
                        value={this.props.service.name}
                        onChange={this.handleChange.bind(this)} />
                </div>
                <div className="form-group">
                    <label htmlFor="endpoint">Endpoint</label>
                    <input
                        className="form-control"
                        type="url"
                        placeholder="url"
                        name="endpoint"
                        ref="endpoint"
                        value={this.props.service.endpoint}
                        onChange={this.handleChange.bind(this)} />
                </div>
                <div className="form-group">
                    <label htmlFor="certificateLocation">Certificate Location</label>
                    <input
                        className="form-control"
                        type="url"
                        placeholder="url"
                        name="certificateLocation"
                        ref="certificateLocation"
                        value={this.props.service.certificateLocation}
                        onChange={this.handleChange.bind(this)} />
                </div>
                <button
                    type="submit"
                    className="btn btn-primary"
                    onClick={this.submit.bind(this)}>submit</button>
                {deleteButton}
           </div>
        );
    }
}

class Service extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            selected: this.defaultCreateData(),
            links: [{name: 'create', id: '1'}],
        };
    }

    defaultCreateData() {
        return {
            id: '1',
            name: '',
            endpoint: '',
            certificateLocation: ''
        };
    }

    componentDidMount() {
        this.updateList();
    }

    updateList() {
        this.props.api.list().then((links) => {
            links.push({name: 'create', id: '1'});
            this.setState({links});
        });
    }

    handleSelection(link) {
        if (link.id === '1') {
            this.setState({
                selected: this.defaultCreateData(),
            });

            return;
        }

        this.props.api.get(link.id).then((service) => {
            this.setState({selected: service || this.defaultCreateData});
        });
    }

    handleChange(service) {
        this.setState({selected: service});
    }

    submit(service) {
        if (service.id == '1') {
            delete service.id;
        }

        this.props.api.save(service).then((service) => {
            const selected = this.state.selected;
            selected.id = service.id;
            this.setState({selected});

            this.updateList();
        });
    }

    delete(service) {
        this.props.api.delete(service).then(() => {
            this.setState({selected: this.defaultCreateData()})
            this.updateList();
        });
    }

    render() {
        const selected = this.state.selected
            ? this.state.selected.id
            : '';

        return (
            <div id="service">
                <Sidebar
                    links={this.state.links}
                    selected={selected}
                    handleSelection={this.handleSelection.bind(this)} />
                <ServiceDetails
                    service={this.state.selected}
                    handleChange={this.handleChange.bind(this)}
                    submit={this.submit.bind(this)}
                    delete={this.delete.bind(this)} />
            </div>
        );
    }
}

class IdpApi {
    constructor() {
        this.counter = 4;
        this.idps = {
            '2': {
                id: '2',
                name: 'test',
                endpoint: 'https://example.com/test/payload',
                certificateLocation: 'https://example.com/test/cert',
            },
            '3': {
                id: '3',
                name: 'test2',
                endpoint: 'https://example.com/test2/payload',
                certificateLocation: 'https://example.com/test2/cert',
            },
        };
    }

    list() {
        return new Promise((resolve) => {
            resolve(Object.keys(this.idps).map((k) => {
                var v = this.idps[k];
                return {id: v.id, name: v.name};
            }));
        });
    }

    get(name) {
        return new Promise((resolve) => {
            resolve(this.idps[name]);
        });
    }

    save(idp) {
        return new Promise((resolve) => {
            if (!idp.id) {
                idp.id = '' + ++this.counter;
            }

            this.idps[idp.id] = idp;
            resolve(idp);
        });
    }

    delete(idp) {
        return new Promise((resolve) => {
            delete this.idps[idp.id];
            resolve();
        });
    }
}

class FederationApi {
    constructor() {
        this.counter = 4;
        this.federations = {
            '2': {
                id: '2',
                name: 'test federation',
                endpoint: 'https://example.com/federation/payload',
                certificateLocation: 'https://example.com/federation/cert',
            },
            '3': {
                id: '3',
                name: 'test2',
                endpoint: 'https://example.com/federation2/payload',
                certificateLocation: 'https://example.com/federation2/cert',
            },
        };
    }

    list() {
        return new Promise((resolve) => {
            resolve(Object.keys(this.federations).map((k) => {
                var v = this.federations[k];
                return {id: v.id, name: v.name};
            }));
        });
    }

    get(name) {
        return new Promise((resolve) => {
            resolve(this.federations[name]);
        });
    }

    save(federation) {
        return new Promise((resolve) => {
            if (!federation.id) {
                federation.id = '' + ++this.counter;
            }

            this.federations[federation.id] = federation;
            resolve(federation);
        });
    }

    delete(federation) {
        return new Promise((resolve) => {
            delete this.federations[federation.id];
            resolve();
        });
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const createIDP = (props) => {
            return <Service api={this.props.idpApi} {...props} />;
        };

        const createFederation = (props) => {
            return <Service api={this.props.federationApi} {...props} />;
        };

        return (
            <Router>
                <div>
                    <nav className="navbar navbar-inverse navbar-fixed-top">
                      <div className="container-fluid">
                        <div className="navbar-header">
                          <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                            <span className="sr-only">Toggle navigation</span>
                            <span className="icon-bar"></span>
                            <span className="icon-bar"></span>
                            <span className="icon-bar"></span>
                          </button>
                          <a className="navbar-brand" href="#">FedProxy</a>
                        </div>
                        <Navbar />
                      </div>
                    </nav>

                    <div className="container-fluid" id="app">
                        <div className="row">
                            <Switch>
                                <Route exact path="/" component={Home} />
                                <Route exact path="/idps" component={createIDP} />
                                <Route exact path="/federations" component={createFederation} />
                            </Switch>
                        </div>
                    </div>
                </div>
            </Router>
        );
    }
}

const app = <App
    idpApi={new IdpApi()}
    federationApi={new FederationApi()} />
ReactDOM.render(app,  document.getElementById('app'))
