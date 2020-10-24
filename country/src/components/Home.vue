<template>
    <div class="hello">
        <!--    头部-->
        <div class="header">
            <div class="header-content">
                <div class="header-img-big">
                    <img class="header-img" src="../assets/header-img.png"/>
                </div>
                <div class="header-house-big">
                    <img class="header-house" src="../assets/header-house.png"/>
                </div>
            </div>
            <div class="header-text">
                <div class="header-text-first">出入境管理系统</div>
                <div class="header-text-second">Immigration management system</div>
            </div>
        </div>
        <!--    导航栏-->
        <div class="nav">
            <div class="nav-content">
                <div @click.stop="goHome()" class="nav-list">
                    <div @click="navline=0" :class="navline==0 ? 'nav-text':'nav-text-second'">申请</div>
                    <div v-show="navline==0" class="nav-line"></div>
                </div>
                <div @click.stop="goDetail()"  class="nav-list">
                    <div @click="navline=2" :class="navline==2 ? 'nav-text':'nav-text-second'">详情</div>
                    <div v-show="navline==2" class="nav-line"></div>
                </div>
                <div class="nav-list">
                    <div @click="navline=1" :class="navline==1 ? 'nav-text':'nav-text-second'">查询</div>
                    <div v-show="navline==1" class="nav-line"></div>
                </div>

                <div class="nav-list">
                    <div @click="goScan()" :class="navline==3 ? 'nav-text':'nav-text-second'">区块链浏览器</div>
                    <div v-show="navline==3" class="nav-line"></div>
                </div>

            </div>
        </div>
        <!--    内容部分-->
        <div class="section">
            <div class="section-content">
                <div class="section-list">
                    <!--            导入私钥-->
                    <div v-if="!account.address" class="section-tab">
                        <div class="section-input"><input v-model="privateKey" placeholder="请输入您的私钥"
                                                          class="section-input-text"/></div>
                        <button v-on:click="initAccount" class="section-btn">导入私钥</button>
                    </div>
                    <div v-if="account.address" class="section-tab">
                        <div class="section-text"><span>地</span>址</div>
                        <div class="section-input">{{account.address}}</div>
                    </div>
                    <!--            身份证号-->
                    <div class="section-tab">
                        <div class="section-text">身份证号</div>
                        <div class="section-input"><input placeholder="请输入您的身份证号" class="section-input-text"/></div>
                    </div>
                    <!--            申请-->
                    <div style="margin-left:-120px">
                        <button @click="actionApply()" class="section-btn-apply">申请</button>
                    </div>

                </div>

            </div>
        </div>
        <div class="loading" v-show="loading==true">
            <img src="../assets/loading.gif"/>
        </div>
    </div>
</template>

<script>
    import Web3 from 'web3'
    import AABI from '../libs/aABI'
    export default {
        name: 'Home',
        props: {
            msg: String
        },
        data() {
            return {
                navline: 0,
                loading: false,
                privateKey: '',
                account: {
                    address: ''
                },
                contract: '0x5221538292372CA706a62EB0D4890d2C85543bE9',
                providerA: 'http://39.98.214.163:8547',
                providerB: 'http://39.98.214.163:8547',
                userInfo: {
                    name: '憨憨',
                    nation: '汉族',
                    number: '1212312312312321',
                    address: '地址地址地址地址地址地址地址地址地址',
                }
            }
        },
        mounted() {
            this.init()
        },
        methods: {
            init() {
                const web3A = new Web3(this.providerA)
                this.web3A = web3A

                const web3B = new Web3(this.providerB)
                this.web3B = web3B

                let masterPrivateKey = localStorage.getItem('masterPrivateKey')
                if (masterPrivateKey) {
                    this.privateKey = masterPrivateKey
                    this.initAccount()
                }
            },
            async initAccount() {
                try {
                    const accountA = this.web3A.eth.accounts.privateKeyToAccount(this.privateKey)
                    this.account = accountA
                    this.web3A.eth.accounts.wallet.add(accountA)


                    const accountB = this.web3A.eth.accounts.privateKeyToAccount(this.privateKey)
                    this.account = accountB
                    this.web3B.eth.accounts.wallet.add(accountB)

                    localStorage.setItem('masterPrivateKey', this.privateKey)
                } catch (error) {
                    this.$Message.error('无效私钥，请重新输入！')
                }
            },
            async actionApply() {
                this.loading = true;
                const contract = new this.web3A.eth.Contract(AABI)
                contract.options.address = this.contract
                console.log(`address`,this.account.address)
                let user = await contract.methods.userInfo(this.account.address).call()
                console.log(`user`,user)
                this.userInfo = user;

                //发起交易
                let data = contract.methods.lock(this.userInfo.id,this.userInfo.nation,this.userInfo.name,this.userInfo.userAddress).encodeABI()
                const tx = {
                    data,
                    from: this.account.address,
                    gas: 200000,
                    value: 0,
                    gasPrice: 1000000,
                    to: this.contract,
                    chainId: 1000
                }
                console.log(`tx`,tx)
                this.loading = true;
                const {transactionHash} = await this.web3A.eth.sendTransaction(tx)
                console.log(`transactionHash`, transactionHash)
                localStorage.setItem('transactionHash', transactionHash)
                this.loading = false;
                this.$router.push({path: '/detail'})
            },
            goHome() {
                this.navline = 0;
                this.$router.push({path: '/home'})
            },
            goDetail() {
                this.navline = 2;
                this.$router.push({path: '/detail'})
            },
            goScan(){
                window.open('http://39.98.50.87:8001/')
            }
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="less">

    .hello {
        position: relative;
        width: 100%;
        min-height: 100vh;
        padding-bottom: 100px;
    }

    .header {
        width: 100%;
        height: 100%;
        position: relative;
        z-index: 50;
        height: 141px;
        background: linear-gradient(354deg, #F7F7DF 0%, #0F8ECE 100%);

        .header-content {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;

            .header-img-big {
                width: 60%;
                height: 141px;

                .header-img {
                    margin-left: 20%;
                    width: 100%;
                    height: 100%;
                }
            }

            .header-house-big {
                width: 30%;
                height: 141px;

                .header-house {
                    margin-left: 10%;
                    width: 100%;
                    height: 100%;
                }
            }
        }

        .header-text {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            position: absolute;
            top: 30%;
            left: 20%;

            .header-text-first {
                font-size: 34px;
                font-family: PingFangSC-Medium, PingFang SC;
                font-weight: 500;
                color: #FFFFFF;
                line-height: 48px;
                letter-spacing: 17px;
            }

            .header-text-second {
                font-size: 19px;
                font-family: NexaBold;
                color: #FFFFFF;
                line-height: 23px;

            }
        }
    }

    .nav {
        z-index: 1000;
        background: linear-gradient(180deg, #1691F9 0%, #0070DC 100%);

        .nav-content {
            padding-left: 20%;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
            padding-top: 15px;
            padding-bottom: 15px;

            .nav-list {
                display: flex;
                flex-direction: column;
                align-items: center;

                .nav-text-second {
                    padding-bottom: 3px;
                    width: 100%;
                    padding-left: 50px;
                    font-size: 18px;
                    font-family: PingFangSC-Regular, PingFang SC;
                    font-weight: 400;
                    color: #FFFFFF;
                    line-height: 25px;
                    cursor: pointer;
                }

                .nav-text {
                    cursor: pointer;
                    padding-bottom: 3px;
                    width: 100%;
                    padding-left: 50px;
                    display: flex;
                    font-size: 18px;
                    font-family: PingFangSC-Semibold, PingFang SC;
                    font-weight: 600;
                    color: #FFFFFF;
                    line-height: 25px;
                }

                .nav-line {
                    //margin-left: 20px;
                    width: 36px;
                    height: 2px;
                    background: #FFFFFF;

                }
            }
        }
    }

    .section {
        background-color: #EFF8FF;

        .section-content {
            padding-left: 20%;
            padding-right: 20%;
            background-color: #fff;

            .section-list {
                display: flex;
                flex-direction: column;
                align-items: center;
                padding-top: 20px;
                padding-left: 115px;

                .section-btn {
                    cursor: pointer;
                    margin-left: 60px;
                    width: 115px;
                    height: 37px;
                    background: linear-gradient(360deg, #0679E4 0%, #16B8FE 100%);
                    border-radius: 23px;
                    color: white;
                    border: 0;
                    outline: none;
                }

                .section-tab {
                    padding-top: 40px;
                    width: 100%;
                    display: flex;
                    flex-direction: row;
                    align-items: center;

                    .section-text {
                        width: 80px;
                        font-size: 16px;
                        font-family: PingFangSC-Semibold, PingFang SC;
                        font-weight: 600;
                        color: #1383FF;
                        line-height: 22px;
                        padding-right: 50px;

                        span {
                            padding-right: 35px
                        }
                    }

                    .section-input {
                        width: 60%;

                        .section-input-text {
                            padding-left: 10px;
                            width: 100%;
                            height: 41px;
                            background: #FFFFFF;
                            box-shadow: 2px 2px 3px 0px rgba(0, 0, 0, 0.11);
                            border-radius: 4px;
                            border: 1px solid #A7D5FF;
                            outline: none;
                        }

                        .section-input-address {
                            padding-top: 30px;
                            padding-left: 10px;
                            width: 100%;
                            height: 66px;
                            background: #FFFFFF;
                            box-shadow: 2px 2px 3px 0px rgba(0, 0, 0, 0.11);
                            border-radius: 4px;
                            border: 1px solid #A7D5FF;

                        }
                    }

                }

                .section-btn-apply {
                    margin-top: 130px;
                    width: 241px;
                    height: 64px;
                    background: linear-gradient(360deg, #FE5560 0%, #FC9055 100%);
                    border-radius: 32px;
                    cursor: pointer;
                    outline: none;
                    border: 0;
                    color: white;
                    font-size: 24px;
                }
            }
        }
    }

    .loading {
        position: absolute;
        top: 0;
        right: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        height: 100%;
        background: #000;
        opacity: 0.1;

        img {
            margin-top: 300px;
        }
    }

</style>
