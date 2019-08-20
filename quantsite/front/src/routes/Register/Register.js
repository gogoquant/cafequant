import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Form } from 'antd';
import classNames from 'classnames';
import Title from 'components/Title';
import { RegisterToken } from './components';
import { register } from './RegisterRedux';
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

  render() {
    return (
      <Fragment>
        <Title title="用户注册" />
        <RegisterToken {...this.props} />
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
