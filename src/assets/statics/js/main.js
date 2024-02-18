function LoadTemplate() {
  var t = Date.now();
  var TEMPLATE_PATH = "/statics/template/";
  var template_list = [
    "main", "list", "show", "toast",
  ];
  return Promise.all(
    _.map(template_list, name => axios.get(TEMPLATE_PATH + name + ".html?t=" + t))
  ).then(function(resps) {
    var tpl = _.reduce(resps, (ret, r, index) => {
      ret[template_list[index]] = r.data;
      return ret
    }, {});
    return tpl;
  });
}

$(function() {
  const BASE_API = "/api/v1";
  function InitStore() {
    const store = new Vuex.Store({
      state: {
        toasts: [],
        toastId: 1,
      },
      mutations: {
        showToast(state, data) {
          state.toastId++;
          state.toasts.push(data);
        },
        hideToast(state, id) {
          state.toasts = _.filter(state.toasts, (o) => o.id !== id);
        },
      },
      actions: {
        showToast(ctx, data) {
          var id = ctx.state.toastId;
          ctx.commit('showToast', {
            id,
            type: data.type,
            message: data.message,
          });
          setTimeout(() => {
            ctx.commit('hideToast', id);
          }, 3000);
        },
      },
    });
    return store;
  }
  function errorToast(store, error, defaultMessage) {
    let msg = _.get(error, 'response.data.Error') || _.get(error, 'message') || defaultMessage
    store.dispatch('showToast', {type: 'error', message: msg});
  }
  LoadTemplate().then((templates) => {
    var Toast = {
      template: templates.toast,
    }
    var MainLayout = {
      template: templates.main,
      components: {
        Toast,
      }
    };
    var ListPage = {
      template: templates.list,
      data: () => {
        return {
          list: [],
          meta: {},
          selectedId: null,
          formMode: 'create',
          formData: {
            Name: '',
            Alias: '',
            Url: '',
            AuthType: '',
            Username: '',
            Password: '',
            SSHKey: '',
          },
          formErrors: {
            Name: '',
            Alias: '',
            Url: '',
            AuthType: '',
            Username: '',
            Password: '',
            SSHKey: '',
          },
        }
      },
      methods: {
        fetchData(page) {
          var meta = this.meta;
          axios.get(`${BASE_API}/repositories/`, {
            params: {
              "pagination[page]": page,
              "pagination[pageSize]": meta.PageSize,
            }
          }).then((res) => {
            this.list = res.data.Data || [];
            this.meta = res.data.Meta.Pagination || {};
          })
        },
        fetchPrev() {
          var meta = this.meta;
          var page = (meta.Page || 1) - 1;
          if (page <= 0) { return }
          this.fetchData(page);
        },
        fetchNext() {
          var meta = this.meta;
          this.fetchData((meta.Page || 1) + 1);
        },
        showDeleteModal(id) {
          document.querySelector("#delete_confirm").showModal();
          this.selectedId = id;
        },
        showEditModal(item) {
          this.formMode = 'update';
          this.selectedId = item.ID;
          this.formData.Name = item.Name;
          this.formData.Alias = item.Alias;
          this.formData.Url = item.Url;
          this.formData.AuthType = item.AuthType === 'sshkey';
          this.formData.Username = item.Username;
          this.formData.Password = item.Password;
          this.formData.SSHKey = item.SSHKey;
          for (const key in this.formErrors) {
            this.formErrors[key] = '';
          }
          document.querySelector("#create_modal").showModal();
        },
        showCreateModal() {
          for (const key in this.formData) {
            this.formData[key] = '';
          }
          this.formData.AuthType = false;
          for (const key in this.formErrors) {
            this.formErrors[key] = '';
          }
          this.formMode = 'create';
          document.querySelector("#create_modal").showModal();
        },
        closeCreateModal() {
          document.querySelector("#create_modal").close();
        },
        async handleDelete() {
          console.log('delete data', this.selectedId);
          try {
            await axios.delete(`${BASE_API}/repositories/${this.selectedId}`);
            this.fetchData(1);
            this.$store.dispatch('showToast', {type: 'success', message: '删除成功'});
          } catch (error) {
            this.$store.dispatch('showToast', {type: 'error', message: '删除失败'});
          }
        },
        async handleSubmitForm() {
          let payload = _.assign({}, this.formData);
          payload.AuthType = this.formData.AuthType ? 'sshkey' : 'password';
          try {
            if (this.formMode === 'update' && this.selectedId) {
              await axios.put(`${BASE_API}/repositories/${this.selectedId}`, payload);
              this.$store.dispatch('showToast', {type: 'success', message: '更新成功'});
            } else {
              await axios.post(`${BASE_API}/repositories/`, payload);
              this.$store.dispatch('showToast', {type: 'success', message: '创建成功'});
            }
            this.closeCreateModal();
            this.fetchData(1);
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
      },
      created: function() {
        console.log("fetch data");
        this.fetchData(1);
      },
    };
    var ShowPage = {
      template: templates.show,
      data: () => {
        return {
          data: {},
          selectedMirrorId: null,
          formMode: 'create',
          formData: {
            Name: '',
            Alias: '',
            Url: '',
            AuthType: '',
            Username: '',
            Password: '',
            SSHKey: '',
          },
          formErrors: {
            Name: '',
            Alias: '',
            Url: '',
            AuthType: '',
            Username: '',
            Password: '',
            SSHKey: '',
          },
          branchList: [],
          tagList: [],
          mirrorList: [],
        }
      },
      methods: {
        async fetchData() {
          try {
            let res = await axios.get(`${BASE_API}/repositories/${this.$route.params.id}`);
            this.data = res.data.Data;
            res = await axios.get(`${BASE_API}/repositories/${this.$route.params.id}/branches/`);
            const refs = _.groupBy(res.data.Data, 'IsTag')
            this.branchList = _.sortBy(refs[false], 'Name')
            this.tagList = _.sortBy(refs[true], 'Name')
            this.fetchMirrors()
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
        async fetchMirrors() {
          try {
            const res = await axios.get(`${BASE_API}/repositories/${this.$route.params.id}/mirrors/`);
            this.mirrorList = res.data.Data;
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
        showDeleteModal(id) {
          document.querySelector("#delete_confirm").showModal();
          this.selectedMirrorId = id;
        },
        showEditModal(item) {
          this.formMode = 'update';
          this.selectedMirrorId = item.ID;
          this.formData.Name = item.Name;
          this.formData.Alias = item.Alias;
          this.formData.Url = item.Url;
          this.formData.AuthType = item.AuthType === 'sshkey';
          this.formData.Username = item.Username;
          this.formData.Password = item.Password;
          this.formData.SSHKey = item.SSHKey;
          for (const key in this.formErrors) {
            this.formErrors[key] = '';
          }
          document.querySelector("#create_modal").showModal();

          this.selectedMirrorId = item.ID;
          for (const key in this.formErrors) {
            this.formErrors[key] = '';
          }
          document.querySelector("#create_modal").showModal();
        },
        showCreateModal() {
          for (const key in this.formData) {
            this.formData[key] = '';
          }
          this.formData.AuthType = false;
          for (const key in this.formErrors) {
            this.formErrors[key] = '';
          }
          this.formMode = 'create';
          document.querySelector("#create_modal").showModal();
        },
        closeCreateModal() {
          document.querySelector("#create_modal").close();
        },
        async handleSubmitForm() {
          let payload = _.assign({}, this.formData);
          payload.Name = payload.Name.trim();
          if (payload.Name === 'origin') {
            this.formErrors.Name = "不能使用 origin";
            return;
          }
          this.formErrors.Name = '';

          payload.AuthType = this.formData.AuthType ? 'sshkey' : 'password';
          try {
            if (this.formMode === 'update' && this.selectedMirrorId) {
              await axios.put(`${BASE_API}/repositories/${this.$route.params.id}/mirrors/${this.selectedMirrorId}`, payload);
              this.$store.dispatch('showToast', {type: 'success', message: '更新成功'});
            } else {
              await axios.post(`${BASE_API}/repositories/${this.$route.params.id}/mirrors/`, payload);
              this.$store.dispatch('showToast', {type: 'success', message: '创建成功'});
            }
            this.closeCreateModal();
            this.fetchMirrors();
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
        async handleDelete() {
          try {
            await axios.delete(`${BASE_API}/repositories/${this.$route.params.id}/mirrors/${this.selectedMirrorId}`);
            this.fetchMirrors();
            this.$store.dispatch('showToast', {type: 'success', message: '删除成功'});
          } catch (error) {
            errorToast(this.$store, error, '删除失败')
          }
        },
        async handlePull() {
          try {
            const res = await axios.put(`${BASE_API}/repositories/${this.$route.params.id}/pull`);
            this.$store.dispatch('showToast', {type: 'info', message: '正在拉取仓库数据'});
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
        async handlePush(item) {
          try {
            const res = await axios.put(`${BASE_API}/repositories/${this.$route.params.id}/mirrors/${item.ID}/push`);
            this.$store.dispatch('showToast', {type: 'info', message: '正在推送仓库数据'});
          } catch (error) {
            errorToast(this.$store, error)
          }
        },
      },
      created: function() {
        console.log("fetch data");
        this.fetchData();
      },
    };

    var router = new VueRouter.createRouter({
      history: VueRouter.createWebHashHistory(),
      routes: [
        {
          path: '/',
          component: MainLayout,
          children: [
            {
              path: '',
              component: ListPage,
              name: 'list',
            },
            {
              path: 'show/:id',
              component: ShowPage,
              name: 'show',
            }
          ],
        }
      ],
      route404: {
        path: '/404',
        component: {
          template: '<div>not found</div>'
        }
      },
    });
    window.store = InitStore();

    window.app = Vue.createApp({});
    app.use(window.store);
    app.use(router);

    app.mount('#app');
  })
});