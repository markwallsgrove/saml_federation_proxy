import React from 'react';
import ReactDOM from 'react-dom';
import axios from 'axios';

function paginateEntityDescriptors() {
    return axios.get('/1/api/entitydescriptor').then(resp => {
        // TODO: error?
        return resp;
    });
}

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            entitydescriptors: [],
        };
    }

    componentDidMount() {
        paginateEntityDescriptors().then(resp => {
            this.setState({entitydescriptors: resp.data});
        });
    }

    render() {
        return (
            <div>
                <h1>Hello, Toast cat!</h1>
                <ul>
                    {this.state.entitydescriptors.map(ed => {
                        return <li>name: {ed.EntityID}</li>
                    })}
                </ul>
            </div>
        );
    }
}

ReactDOM.render(<App />,  document.getElementById('app'))
