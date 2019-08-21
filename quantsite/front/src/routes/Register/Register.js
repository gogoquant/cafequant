import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Form } from 'antd';
import classNames from 'classnames';
import Title from 'components/Title';
import { RegisterToken } from './components';
import { register } from './RegisterRedux';
import { UserService } from 'services';
import styles from './Register.scss';

const cls = classNames(styles.container, styles.xsContainer);

export class RegisterComponent extends React.Component {
  static childContextTypes = {
    form: PropTypes.object,
  };

  constructor(props) {
    super(props);
  }

  getChildContext() {
    return {
      form: this.props.form,
    };
  }

  handleSubmit = e => {
    e.preventDefault();
    this.props.form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values);
        //try to call register
        UserService.createUser(values.email, values.password, values.nickname);
        //this.setState({ values: values });
      }
    });
  };

  render() {
    return (
      <Fragment>
        <Title title="用户注册" />
        <RegisterToken {...this.props} handleSubmit={this.handleSubmit} />
      </Fragment>
    );
  }
}

RegisterComponent.propTypes = {
  history: PropTypes.shape({
    push: PropTypes.func,
  }),
  form: PropTypes.object,
  loading: PropTypes.bool,
  userRegister: PropTypes.func,
};

export default withRouter(
  Form.create()(
    connect(
      state => ({}),
      {
        userRegister: register,
      },
    )(RegisterComponent),
  ),
);
