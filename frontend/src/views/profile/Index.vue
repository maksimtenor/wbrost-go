<script setup>

import Navbar from "../../components/layout/Navbar.vue";
import Sidebar from "../../components/layout/Sidebar.vue";
</script>

<template>
  <Navbar />
  <!-- Main Sidebar Container -->
  <Sidebar />
  <div class="content-wrapper">
  <div class="content-header">
    <div class="container-fluid">
      <div class="row mb-2">
        <div class="col-sm-6">
          <h1 class="m-0">
            Личный кабинет                    </h1>
        </div><!-- /.col -->
        <div class="col-sm-6">
        </div><!-- /.col -->
      </div><!-- /.row -->
    </div><!-- /.container-fluid -->
  </div>
  <div class="content">



    <form id="profile-form" action="/profile/index" method="post">
      <input type="hidden" name="_csrf" value="SdeOyQg4nFz-YhQ9suokgkbzOt_j8-NxFeyDVuFRHX4ep6OzW1bfA58rJHeAuXHoDLt3sbagjCNzirEf0hpxPQ==">
      <div class="form-group field-user-name">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-name">Имя:</label>
        <input type="text" id="user-name" class="col-lg-3 form-control" name="User[name]" value="System" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group field-user-username">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-username">Логин:</label>
        <input type="text" id="user-username" class="col-lg-3 form-control" name="User[username]" value="admin" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group field-user-password">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-password">Смена пароля:</label>
        <input type="password" id="user-password" class="col-lg-3 form-control" name="User[password]" value="" placeholder="Введите новый" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group field-user-email">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-email">Email:</label>
        <input type="text" id="user-email" class="col-lg-3 form-control" name="User[email]" value="maksimrozette@gmail.com" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group field-user-phone">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-phone">Номер телефона:</label>
        <input type="text" id="user-phone" class="col-lg-3 form-control" name="User[phone]" value="+79937799526" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="mb-3">
        <label class="form-label">PRO-аккаунт:</label>
        <span class="form-control" style="color: green;width: 400px">Активен</span>    </div>

      <div class="form-group field-user-taxes">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-taxes">Налоговый %:</label>
        <input type="text" id="user-taxes" class="col-lg-3 form-control" name="User[taxes]" value="5" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group field-user-wb_key">
        <label class="col-lg-3 col-form-label mr-lg-3" for="user-wb_key">WB API KEY:</label>
        <input type="text" id="user-wb_key" class="col-lg-3 form-control" name="User[wb_key]" value="eyJhbGciOiJFUzI1NiIsImtpZCI6IjIwMjUwNTIwdjEiLCJ0eXAiOiJKV1QifQ.eyJlbnQiOjEsImV4cCI6MTc3MDk4NTk0MSwiaWQiOiIwMTk4YWIyNC01NTg3LTdhZGEtYmM0Yi1mZTc0NTA0YmM5ZmYiLCJpaWQiOjg1MzY3NjYxLCJvaWQiOjEzMjE3ODQsInMiOjE2MTI2LCJzaWQiOiJiZjE0NjQzMC04NmVmLTRlNDAtYmNmNS1lNDY1Yzk1MjU4ZWEiLCJ0IjpmYWxzZSwidWlkIjo4NTM2NzY2MX0.S4mgg6DGXzHp7FqYIBWKpLJa2O-WYNan1tdQf6pQy8S8HcOZBIDFc_q3S-tuWuUA8pIT3amtyfgiK2YGdYfiPg" placeholder="Введите" style="width: 400px">
        <p class="col-lg-7 invalid-feedback"></p>
      </div>
      <div class="form-group">
        <div>
          <button type="submit" class="btn btn-primary form-button">Сохранить</button>    </div>
      </div>

    </form><!-- /.container-fluid -->
  </div>
  </div>
</template>

<style scoped>

</style>