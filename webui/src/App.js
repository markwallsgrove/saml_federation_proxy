import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';

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

function toggleEntityDescriptor(id, enabled) {
    const url = '/1/api/entityDescriptor/toggle' +
        '?enabled=' + enabled +
        '&id=' + id;

    return axios.post(url).then(resp => {
        // TODO: error?
        return resp.data;
    });
}

class EntityDescriptor extends React.Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    toggleEntityDescriptor() {
        toggleEntityDescriptor(this.props.id, !this.props.enabled).then(ed => {
            this.props.entityChanged(ed);
        });
    }

    render() {
        const entityChanged = this.toggleEntityDescriptor.bind(this);
        return (<li onClick={entityChanged}>{this.props.name}(enabled={this.props.enabled})</li>);
    }
}

class EntityDescriptors extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            entitydescriptors: [],
        };
    }

    componentDidMount() {
        paginateEntityDescriptors().then(resp => {
            this.setState({entitydescriptors: resp});
        });
    }

    entityChanged(entityDescriptor) {
        console.log("entity changed ", entityDescriptor);
        // TODO: change entitydescriptors structor to use key/value map
    }

    render() {
        return (
            <ul>
                {this.state.entitydescriptors.map(ed => {
                    console.log('ed', ed.Enabled);
                    return (
                        <EntityDescriptor
                            federationID={ed.FederationID}
                            id={ed.EntityID}
                            name={ed.EntityID}
                            enabled={ed.Enabled}
                            entityChanged={this.entityChanged.bind(this)}
                        />
                    );
                })}
            </ul>
        );
    }
}

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <h1>Hello, Toast cat!</h1>
                <EntityDescriptors />
            </div>
        );
    }
}

ReactDOM.render(<App />,  document.getElementById('app'))
